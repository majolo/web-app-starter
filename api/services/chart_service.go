package services

import (
	"bytes"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/csv"
	"errors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/majolo/growth-charts/generated/go/charts"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"regexp"
	"strings"
	"time"
)

type Template struct {
	pages       []Page    `firestore:"pages,omitempty"`
	variables   []string  `firestore:"variables,omitempty"`
	name        string    `firestore:"name"`
	createdAt   time.Time `firestore:"time,omitempty"`
	primaryFont string    `firestore:"primary_font,omitempty"`
}

type Page struct {
	Content string `firestore:"content"`
	Emoji   string `firestore:"emoji"`
	// Colour can be ROYGBIV
	Colour string `firestore:"colour,omitempty"`
}

type ChartService struct {
	webDomain string
	charts.UnimplementedChartsServer
	fsClient *firestore.Client
}

const (
	templateCollection    = "templates"
	instanceSubCollection = "instances"
)

func NewChartService(fsClient *firestore.Client, webDomain string) *ChartService {
	// Pages for example template
	_ = []Page{
		{
			Content: "${name}'s 2022 Alphabeticle Recap!",
			Emoji:   "🅰️",
			Colour:  "Red",
		},
		{
			Content: "You played a total of ${games_played} games.",
			Emoji:   "📈",
			Colour:  "Violet",
		},
		{
			Content: "You won ${wins}, losing only ${losses}!",
			Emoji:   "🤑",
			Colour:  "Green",
		},
		{
			Content: "Mr Worldwide! You remembered to play from ${countries} countries.",
			Emoji:   "🌎",
			Colour:  "Yellow",
		},
		{
			Content: "${starter_word} was your favourite starter word.",
			Emoji:   "🏃",
			Colour:  "Orange",
		},
		{
			Content: "The least common word played by you was ${rarest_word}.",
			Emoji:   "📔",
			Colour:  "Blue",
		},
		{
			Content: "Your favourite day to play was ${favourite_day}.",
			Emoji:   "📅",
			Colour:  "Indigo",
		},
		{
			Content: "Longest streak of ${longest_streak} days, better than ${percentile}% of players.",
			Emoji:   "🏋",
			Colour:  "Red",
		},
		{
			Content: "You unlocked ${trophies} trophies!",
			Emoji:   "🏆",
			Colour:  "Orange",
		},
		{
			Content: "How will you do in 2023?",
			Emoji:   "😀",
			Colour:  "Yellow",
		},
	}

	cs := &ChartService{
		fsClient:  fsClient,
		webDomain: webDomain,
	}
	//_, _ = cs.CreateTemplate(context.Background(), &charts.CreateTemplateRequest{
	//	Name:  "Alphabeticle test v1",
	//	Pages: nativePageToProtoPage(pages),
	//})
	return cs
}

func (c *ChartService) GetInstance(ctx context.Context, in *charts.GetInstanceRequest) (*charts.GetInstanceResponse, error) {
	if in.InstanceId == "test" {
		resp := charts.GetInstanceResponse{
			Variables: map[string]string{},
		}
		return &resp, nil
	}

	resp := charts.GetInstanceResponse{
		Variables: nil,
	}

	d, err := c.fsClient.Collection(templateCollection).Doc(in.TemplateId).Collection(instanceSubCollection).Doc(in.InstanceId).Get(ctx)
	if err != nil {
		return nil, err
	}
	vars := map[string]string{}
	for k, v := range d.Data() {
		vars[k] = v.(string)
	}
	resp.Variables = vars

	return &resp, nil
}

func (c *ChartService) CreateInstance(ctx context.Context, in *charts.CreateInstanceRequest) (*charts.CreateInstanceResponse, error) {
	resp := charts.CreateInstanceResponse{
		InstanceId: "",
	}
	d, _, err := c.fsClient.Collection(templateCollection).Doc(in.TemplateId).Collection(instanceSubCollection).Add(ctx, in.Variables)
	if err != nil {
		return nil, err
	}
	resp.InstanceId = d.ID
	return &resp, nil
}

func (c *ChartService) GenerateTemplateCSV(ctx context.Context, in *charts.GenerateTemplateCSVRequest) (*charts.GenerateTemplateCSVResponse, error) {
	resp := charts.GenerateTemplateCSVResponse{}
	tmp, err := c.GetTemplate(ctx, &charts.GetTemplateRequest{TemplateId: in.TemplateId})
	if err != nil {
		return nil, err
	}
	resp.Csv = strings.Join(tmp.Template.Variables, ",")
	return &resp, nil
}

func (c *ChartService) CreateInstancesFromCSV(ctx context.Context, in *charts.CreateInstancesFromCSVRequest) (*charts.CreateInstancesFromCSVResponse, error) {
	resp := charts.CreateInstancesFromCSVResponse{
		NewCsv: "",
	}
	//todo: verify that the csv titles are valid

	newIds := []string{}
	var firstRow []string

	csvReader := csv.NewReader(strings.NewReader(in.Csv))
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(firstRow) == 0 {
			firstRow = rec
			continue
		}

		vars := map[string]string{}
		for i, v := range rec {
			vars[firstRow[i]] = v
		}

		resp, err := c.CreateInstance(ctx, &charts.CreateInstanceRequest{
			TemplateId: in.TemplateId,
			Variables:  vars,
		})
		if err != nil {
			return nil, err
		}
		newIds = append(newIds, resp.InstanceId)
	}

	b := new(bytes.Buffer)
	w := csv.NewWriter(b)

	csvReader = csv.NewReader(strings.NewReader(in.Csv))
	firstRow = []string{}
	i := 0

	pairs := [][]string{}

	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(firstRow) == 0 {
			firstRow = rec
			firstRow = append(firstRow, "url")
			pairs = append(pairs, firstRow)
			continue
		}

		var temp []string
		temp = rec
		uri := c.webDomain + "/render/" + in.TemplateId + "/" + newIds[i]
		temp = append(temp, uri)
		pairs = append(pairs, temp)
		i++
	}
	err := w.WriteAll(pairs)
	if err != nil {
		return nil, err
	}
	resp.NewCsv = b.String()
	return &resp, nil
}

func (c *ChartService) CreateTemplate(ctx context.Context, in *charts.CreateTemplateRequest) (*charts.CreateTemplateResponse, error) {
	resp := charts.CreateTemplateResponse{
		TemplateId: "",
	}
	vars := extractVars(protoPageToNativePage(in.Pages))
	// TODO: for some reason firestore doesn't like this, would be better but can fix at a later date
	//tmpl := Template{
	//	pages:     protoPageToNativePage(in.Pages),
	//	variables: vars,
	//	name:      in.Name,
	//	createdAt: time.Now(),
	//}
	d, _, err := c.fsClient.Collection(templateCollection).Add(ctx, map[string]interface{}{
		"pages":       protoPageToNativePage(in.Pages),
		"variables":   vars,
		"name":        in.Name,
		"createdAt":   time.Now(),
		"primaryFont": in.PrimaryFont,
	})
	if err != nil {
		return nil, err
	}
	resp.TemplateId = d.ID
	return &resp, nil
}

func extractVars(pages []Page) []string {
	vars := []string{}
	for _, p := range pages {
		// regex for checking all ${x} in p.Content
		re := regexp.MustCompile(`\${(.*?)}`)
		matches := re.FindAllStringSubmatch(p.Content, -1)
		for _, m := range matches {
			vars = append(vars, m[1])
		}
	}
	// remove duplicates
	vars = unique(vars)
	return vars
}

func unique(s []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range s {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func (c *ChartService) GetTemplate(ctx context.Context, in *charts.GetTemplateRequest) (*charts.GetTemplateResponse, error) {
	resp := charts.GetTemplateResponse{
		Template: nil,
	}
	doc, err := c.fsClient.Collection(templateCollection).Doc(in.TemplateId).Get(ctx)
	if err != nil {
		return nil, err
	}
	if !stringMapContainsKeys(doc.Data(), []string{"name", "createdAt", "pages", "variables"}) {
		return nil, errors.New("invalid template")
	}
	// TODO: refactor this out with list also
	pages := []Page{}
	for _, v := range doc.Data()["pages"].([]interface{}) {
		pages = append(pages, Page{
			Colour:  v.(map[string]interface{})["colour"].(string),
			Content: v.(map[string]interface{})["content"].(string),
			Emoji:   v.(map[string]interface{})["emoji"].(string),
		})
	}
	vars := []string{}
	for _, v := range doc.Data()["variables"].([]interface{}) {
		vars = append(vars, v.(string))
	}
	var primaryFont string
	if val, ok := doc.Data()["primaryFont"]; ok {
		primaryFont = val.(string)
	} else {
		primaryFont = "Segoe UI"
	}

	resp.Template = &charts.Template{
		TemplateId:  doc.Ref.ID,
		Name:        doc.Data()["name"].(string),
		CreatedAt:   timestamppb.New(doc.Data()["createdAt"].(time.Time)),
		Pages:       nativePageToProtoPage(pages),
		Variables:   vars,
		PrimaryFont: primaryFont,
	}
	return &resp, nil
}

func (c *ChartService) ListTemplates(ctx context.Context, _ *charts.ListTemplatesRequest) (*charts.ListTemplatesResponse, error) {
	resp := charts.ListTemplatesResponse{
		Templates: nil,
	}
	iter := c.fsClient.Collection(templateCollection).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		if !stringMapContainsKeys(doc.Data(), []string{"name", "createdAt", "pages", "variables"}) {
			continue
		}
		pages := []Page{}
		for _, v := range doc.Data()["pages"].([]interface{}) {
			// TODO: what in the name of fuck is this
			pages = append(pages, Page{
				Colour:  v.(map[string]interface{})["colour"].(string),
				Content: v.(map[string]interface{})["content"].(string),
				Emoji:   v.(map[string]interface{})["emoji"].(string),
			})
		}
		vars := []string{}
		for _, v := range doc.Data()["variables"].([]interface{}) {
			vars = append(vars, v.(string))
		}
		resp.Templates = append(resp.Templates, &charts.Template{
			TemplateId: doc.Ref.ID,
			Name:       doc.Data()["name"].(string),
			CreatedAt:  timestamppb.New(doc.Data()["createdAt"].(time.Time)),
			Pages:      nativePageToProtoPage(pages),
			Variables:  vars,
		})
	}
	return &resp, nil
}

func (c *ChartService) DeleteInstance(ctx context.Context, in *charts.DeleteInstanceRequest) (*charts.DeleteInstanceResponse, error) {
	resp := charts.DeleteInstanceResponse{}
	_, err := c.fsClient.Collection(instanceSubCollection).Doc(in.InstanceId).Delete(ctx)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
func (c *ChartService) DeleteTemplate(ctx context.Context, in *charts.DeleteTemplateRequest) (*charts.DeleteTemplateResponse, error) {
	resp := charts.DeleteTemplateResponse{}
	_, err := c.fsClient.Collection(templateCollection).Doc(in.TemplateId).Delete(ctx)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *ChartService) UpdateTemplate(ctx context.Context, in *charts.UpdateTemplateRequest) (*charts.UpdateTemplateResponse, error) {
	resp := charts.UpdateTemplateResponse{}
	pages := protoPageToNativePage(in.Template.Pages)
	_, err := c.fsClient.Collection(templateCollection).Doc(in.TemplateId).Set(ctx, map[string]interface{}{
		"pages":     pages,
		"variables": extractVars(pages),
		"name":      in.Template.Name,
	})
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func helperCompareListsAsSets(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	bAsSet := make(map[string]bool)
	for _, v := range b {
		bAsSet[v] = true
	}

	for _, v := range a {
		if _, ok := bAsSet[v]; !ok {
			return false
		}
	}
	return true
}

func (c *ChartService) RegisterGRPC(gs *grpc.Server) {
	charts.RegisterChartsServer(gs, c)
}

func (c *ChartService) RegisterGRPCGateway(ctx context.Context, mux *runtime.ServeMux) {
	charts.RegisterChartsHandlerServer(ctx, mux, c)
}

func protoPageToNativePage(protoPages []*charts.Page) []Page {
	var pages []Page
	for _, p := range protoPages {
		pages = append(pages, Page{
			Content: p.Content,
			Emoji:   p.Emoji,
			Colour:  p.Colour,
		})
	}
	return pages
}

func nativePageToProtoPage(nativePages []Page) []*charts.Page {
	var pages []*charts.Page
	for _, p := range nativePages {
		pages = append(pages, &charts.Page{
			Content: p.Content,
			Emoji:   p.Emoji,
			Colour:  p.Colour,
		})
	}
	return pages
}

func stringMapContainsKeys(m map[string]interface{}, keys []string) bool {
	for _, k := range keys {
		if _, ok := m[k]; !ok {
			return false
		}
	}
	return true
}
