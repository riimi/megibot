package template

import (
	"fmt"
	"github.com/line/line-bot-sdk-go/linebot"
	"kitchensink/model"
	"log"
)

func BubbleContainerRssService(service model.RssService, nowsub bool) *linebot.BubbleContainer {
	ret := &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		//		Hero: &linebot.ImageComponent{
		//			Type:        linebot.FlexComponentTypeImage,
		//			URL:         service.Url,
		//			Size:        linebot.FlexImageSizeTypeFull,
		//			AspectRatio: linebot.FlexImageAspectRatioType20to13,
		//			AspectMode:  linebot.FlexImageAspectModeTypeCover,
		//		},
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   service.Name,
					Size:   linebot.FlexTextSizeTypeXl,
					Weight: linebot.FlexTextWeightTypeBold,
				},
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeVertical,
					Contents: []linebot.FlexComponent{
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeBaseline,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  "Title",
									Flex:  1,
									Size:  linebot.FlexTextSizeTypeSm,
									Color: "#aaaaaa",
								},
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  service.Title,
									Flex:  5,
									Size:  linebot.FlexTextSizeTypeSm,
									Wrap:  true,
									Color: "#666666",
								},
							},
							Spacing: linebot.FlexComponentSpacingTypeSm,
						},
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeBaseline,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  "Desc",
									Flex:  1,
									Size:  linebot.FlexTextSizeTypeSm,
									Color: "#aaaaaa",
								},
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  service.Description,
									Flex:  5,
									Size:  linebot.FlexTextSizeTypeSm,
									Wrap:  true,
									Color: "#666666",
								},
							},
							Spacing: linebot.FlexComponentSpacingTypeSm,
						},
					},
					Spacing: linebot.FlexComponentSpacingTypeSm,
					Margin:  linebot.FlexComponentMarginTypeLg,
				},
			},
		},
		Footer: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.ButtonComponent{
					Type: linebot.FlexComponentTypeButton,
					Action: &linebot.PostbackAction{
						Label:       "Subscribe",
						Data:        fmt.Sprintf(`Rss subscribe %s`, service.Name),
						DisplayText: fmt.Sprintf(`Rss subscribe %s`, service.Name),
					},
					Height: linebot.FlexButtonHeightTypeSm,
					Style:  linebot.FlexButtonStyleTypePrimary,
				},
				&linebot.ButtonComponent{
					Type: linebot.FlexComponentTypeButton,
					Action: &linebot.URIAction{
						Label: "Website",
						URI:   service.Url,
					},
					Height: linebot.FlexButtonHeightTypeSm,
					Style:  linebot.FlexButtonStyleTypeSecondary,
				},
				&linebot.SpacerComponent{
					Type: linebot.FlexComponentTypeSpacer,
					Size: linebot.FlexSpacerSizeTypeSm,
				},
			},
			Spacing: linebot.FlexComponentSpacingTypeSm,
		},
	}
	if nowsub {
		ret.Footer.Contents[0] = &linebot.ButtonComponent{
			Type: linebot.FlexComponentTypeButton,
			Action: &linebot.PostbackAction{
				Label:       "Unsubscribe",
				Data:        fmt.Sprintf(`Rss unsubscribe %s`, service.Name),
				DisplayText: fmt.Sprintf(`Rss unsubscribe %s`, service.Name),
			},
			Height: linebot.FlexButtonHeightTypeSm,
			Style:  linebot.FlexButtonStyleTypeSecondary,
		}
	}

	return ret
}

func FlexContainerRssServices(services []model.RssService, id string) linebot.FlexContainer {
	subs, err := model.AllSubsByUser(id)
	if err != nil {
		log.Print(err)
	}
	subsMap := make(map[string]bool)
	for _, sub := range subs {
		subsMap[sub.ServiceName] = true
	}

	contents := make([]*linebot.BubbleContainer, 0, len(services))
	for _, s := range services {
		_, nowsub := subsMap[s.Name]
		contents = append(contents, BubbleContainerRssService(s, nowsub))
	}
	return &linebot.CarouselContainer{
		Type:     linebot.FlexContainerTypeCarousel,
		Contents: contents,
	}
}

func BubbleContainerRssItem(item model.RssItem) *linebot.BubbleContainer {
	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeHorizontal,
			Contents: []linebot.FlexComponent{
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeVertical,
					Flex:   4,
					Contents: []linebot.FlexComponent{
						&linebot.TextComponent{
							Type:   linebot.FlexComponentTypeText,
							Text:   item.ServiceName,
							Size:   linebot.FlexTextSizeTypeSm,
							Weight: linebot.FlexTextWeightTypeBold,
							Color:  "#1DB446",
						},
						&linebot.TextComponent{
							Type:   linebot.FlexComponentTypeText,
							Text:   item.Title,
							Size:   linebot.FlexTextSizeTypeSm,
							Weight: linebot.FlexTextWeightTypeBold,
							Wrap:   true,
						},
						&linebot.TextComponent{
							Type:  linebot.FlexComponentTypeText,
							Text:  item.Published,
							Size:  linebot.FlexTextSizeTypeXs,
							Color: "#aaaaaa",
							Wrap:  true,
						},
					},
				},
				&linebot.ButtonComponent{
					Type: linebot.FlexComponentTypeButton,
					Action: &linebot.URIAction{
						Label: "Web",
						URI:   item.Link,
					},
					Style:   linebot.FlexButtonStyleTypeSecondary,
					Flex:    1,
					Gravity: linebot.FlexComponentGravityTypeCenter,
					Margin:  linebot.FlexComponentMarginTypeSm,
				},
			},
		},
	}
}

/*
func BubbleContainerRssItem(item model.RssItem) *linebot.BubbleContainer {
	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   item.ServiceName,
					Size:   linebot.FlexTextSizeTypeSm,
					Weight: linebot.FlexTextWeightTypeBold,
					Color:  "#1DB446",
				},
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   item.Title,
					Size:   linebot.FlexTextSizeTypeMd,
					Weight: linebot.FlexTextWeightTypeBold,
					Wrap:   true,
				},
				&linebot.TextComponent{
					Type:  linebot.FlexComponentTypeText,
					Text:  item.Published,
					Size:  linebot.FlexTextSizeTypeXs,
					Color: "#aaaaaa",
					Wrap:  true,
				},
			},
		},
		Footer: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.ButtonComponent{
					Type: linebot.FlexComponentTypeButton,
					Action: &linebot.URIAction{
						Label: "Website",
						URI:   item.Link,
					},
					Height: linebot.FlexButtonHeightTypeSm,
					Style:  linebot.FlexButtonStyleTypeLink,
				},
				&linebot.SpacerComponent{
					Type: linebot.FlexComponentTypeSpacer,
					Size: linebot.FlexSpacerSizeTypeSm,
				},
			},
			Spacing: linebot.FlexComponentSpacingTypeSm,
		},
	}
}
*/
