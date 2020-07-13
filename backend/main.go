package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber"
	"golang.org/x/net/html"
)

type Hackathon struct {
	Link string `json:"link"`

	Name string `json:"name"`
}
type SubmissionPeriod struct {
	Begins string `json:"begins"`
	Ends   string `json:"ends"`
}

type Flag struct {
	Severity string `json:severity`
	Message  string `json:message`
}
type OverallRating struct {
	Flags  []Flag `json:flags`
	Rating int32  `json:rating`
}

func main() {

	app := fiber.New()
	//THIS RETURNS ALL COMPxETITIONS IT WAS SUBMITTED TO
	app.Get("/competitions/:software", func(c *fiber.Ctx) {
		fmt.Println("CALLED US")

		var software string = c.Params("software")
		var sofwareURL string = "https://devpost.com/software/" + software
		resp, err := http.Get(sofwareURL)
		if err != nil {
			fmt.Println("FAILED TO CALL: " + sofwareURL)
			fmt.Println(err)
			c.SendStatus(418)

		}
		data, _ := ioutil.ReadAll(resp.Body)
		s := string(data)

		resp.Body.Close()

		dataReader := strings.NewReader(s)

		doc, err := html.Parse(dataReader)
		if err != nil {
			fmt.Println("FAILED TO GET HTML: " + sofwareURL)

			c.SendStatus(418)
		}

		//software-header
		if elementWithIDExists(doc, "software-header") {
			var element *html.Node = getElementById(doc, "submissions")
			var nodesWithin []*html.Node = buildLinkNodes(element)
			var hackathons []Hackathon
			for _, n := range nodesWithin {
				hackathons = append(hackathons, buildHackathon(n))
			}

			hackathons = cleanHackathons(hackathons)
			fmt.Println(hackathons)
			var hackathonJSON []map[string]string = hackathonsToJSON(hackathons)
			if err != nil {
				log.Fatal("Cannot encode to JSON ", err)
			}
			// hackathonJSONStr := string(hackathonJSON)
			fmt.Println(hackathonJSON)
			returnValue, err := json.Marshal(hackathonJSON)
			if err != nil {
				log.Fatal("Cannot encode to JSON ", err)
			}
			c.Send(returnValue)
		} else {
			//WRONG URL

			c.SendStatus(418)

		}

	})

	app.Get("/check/:software", func(c *fiber.Ctx) {
		fmt.Println("Chheck was called")
		var software string = c.Params("software", "")
		if software == "" {
			c.SendStatus(418)
		}

		var hackathon string = c.Query("hackathon", "")
		if hackathon == "" {
			c.SendStatus(418)
		}

		var sofwareURL string = "https://devpost.com/software/" + software
		var hackathonURL string = "https://" + hackathon + ".devpost.com/details/dates"

		softwareRes, softwareErr := http.Get(sofwareURL)
		hackathonRes, hackathonErr := http.Get(hackathonURL)
		if softwareErr != nil || hackathonErr != nil {
			fmt.Println("FAILED TO CALL: " + sofwareURL)
			fmt.Println(softwareErr)
			c.SendStatus(418)

		}
		hackathonData, _ := ioutil.ReadAll(hackathonRes.Body)
		softwareData, _ := ioutil.ReadAll(softwareRes.Body)

		hackathonDataString := string(hackathonData)
		softwareDataString := string(softwareData)

		hackathonRes.Body.Close()
		softwareRes.Body.Close()

		hackathonDataReader := strings.NewReader(hackathonDataString)
		softwareDataReader := strings.NewReader(softwareDataString)

		hackathonDoc, hackathonErr := html.Parse(hackathonDataReader)
		softwareDoc, softwareErr := html.Parse(softwareDataReader)

		if softwareErr != nil && hackathonErr != nil {
			fmt.Println("FAILED TO GET HTML: " + sofwareURL)

			c.SendStatus(418)
		}

		if elementWithIDExists(softwareDoc, "software-header") && elementWithIDExists(hackathonDoc, "main") {
			var submissionPeriod SubmissionPeriod = getSubmssionPeriod(getElementById(hackathonDoc, "main"))
			//c.Send("ITS RUDE NOT TO SEND A RESPONSE")

			parsedStartTime, startTimeErr := time.Parse("January 02 at 3:04am MST", submissionPeriod.Begins)
			parsedEndTime, endTimeErr := time.Parse("January 02 at 3:04am MST", submissionPeriod.Begins)

			if startTimeErr != nil {
				fmt.Println(startTimeErr)
				log.Fatal("START TIME ERROR")
			}
			if endTimeErr != nil {
				fmt.Println(endTimeErr)
				log.Fatal("END TIME ERROR")
			}
			parsedStartTime = parsedStartTime.AddDate(2020, 0, 0)
			parsedEndTime = parsedEndTime.AddDate(2020, 0, 0)

			var flags []Flag
			var rating int32

			containsGithubRepo := containsGithubRepo(softwareDoc)
			if !containsGithubRepo {
				var overallRating OverallRating

				var newFlag Flag
				newFlag.Message = "No Github Repo Listed. This is a severe issue because they provided no verification that their code exists."
				newFlag.Severity = "severe"

				flags = append(flags, newFlag)

				overallRating.Flags = flags

				rating = 3

				overallRating.Rating = rating
				returnValue, err := json.Marshal(overallRating)
				if err != nil {
					log.Fatal("Cannot encode to JSON ", err)
				}
				c.Send(returnValue)
			} else {
				var overallRating OverallRating

				overallRating.Flags = flags

				rating = 0

				overallRating.Rating = rating
				returnValue, err := json.Marshal(overallRating)
				if err != nil {
					log.Fatal("Cannot encode to JSON ", err)
				}
				c.Send(returnValue)
			}

		} else {
			//WRONG URL

			c.SendStatus(418)

		}
	})

	app.Listen(3000)
}

func GetAttribute(n *html.Node, key string) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}
	return "", false
}

func checkId(n *html.Node, id string) bool {
	if n.Type == html.ElementNode {
		s, ok := GetAttribute(n, "id")
		if ok && s == id {
			return true
		}
	}
	return false
}

func traverse(n *html.Node, id string) *html.Node {
	if checkId(n, id) {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result := traverse(c, id)
		if result != nil {
			return result
		}
	}

	return nil
}

func getElementById(n *html.Node, id string) *html.Node {
	return traverse(n, id)
}

func elementWithIDExists(n *html.Node, id string) bool {
	return getElementById(n, id) != nil
}
func buildNodes(n *html.Node, nodeType string) []*html.Node {
	if n.Type == html.ElementNode && n.Data == nodeType && !attrContains(n.Attr, "src") {
		return []*html.Node{n}
	}

	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, buildNodes(c, nodeType)...)
	}
	return ret
}
func buildLinkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" && !attrContains(n.Attr, "src") {
		return []*html.Node{n}
	}

	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, buildLinkNodes(c)...)
	}
	return ret
}

func buildTRNodes(n *html.Node) []*html.Node {
	return buildNodes(n, "tr")
}
func buildHackathon(n *html.Node) (hackathon Hackathon) {
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			hackathon.Link = attr.Val

		}
	}

	hackathon.Name = n.FirstChild.Data

	return
}

func attrContains(s []html.Attribute, e string) bool {
	for _, a := range s {
		if a.Key == e {
			return true
		}
	}
	return false
}
func cleanHackathons(hackathons []Hackathon) []Hackathon {
	var cleanHacks []Hackathon

	for _, n := range hackathons {
		if n.Name != "img" {

			cleanHacks = append(cleanHacks, n)
		}
	}
	return cleanHacks
}

func hackathonsToJSON(hackathons []Hackathon) []map[string]string {
	var jsonHacks []map[string]string

	for _, n := range hackathons {
		hackathonMap := make(map[string]string)
		hackathonMap["link"] = n.Link
		hackathonMap["name"] = n.Name
		jsonHacks = append(jsonHacks, hackathonMap)
	}
	return jsonHacks

}

func getSubmssionPeriod(n *html.Node) SubmissionPeriod {
	var submissionPeriod SubmissionPeriod
	fmt.Println("NODE CHECKIN")
	// fmt.Println(n.Data)
	tableRows := buildTRNodes(n)
	submissionRow := tableRows[1]
	submissionRowTDs := buildNodes(submissionRow, "td")
	submissionPeriod.Begins = submissionRowTDs[1].FirstChild.Data
	submissionPeriod.Ends = submissionRowTDs[2].FirstChild.Data
	return submissionPeriod

}

func containsGithubRepo(n *html.Node) bool {
	links := buildLinkNodes(n)
	for i := 0; i < len(links); i++ {
		link := links[i]
		if attrContains(link.Attr, "href") {
			href, _ := GetAttribute(link, "href")
			if strings.Contains(href, "https://github.com") {
				return true
			}
		}

	}
	return false
}
