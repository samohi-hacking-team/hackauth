package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/html"

	"github.com/gofiber/fiber"
)

func main() {

	app := fiber.New()
	//THIS RETURNS ALL COMPxETITIONS IT WAS SUBMITTED TO
	app.Get("/check/:software", func(c *fiber.Ctx) {
		var software string = c.Params("software")
		var sofwareURL string = "https://devpost.com/software/" + software
		resp, err := http.Get(sofwareURL)
		if err != nil {
			fmt.Println("FAILED TO CALL: " + sofwareURL)
			fmt.Println(err)
			c.SendStatus(418)

		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {

			fmt.Println("FAILED TO PARSE: " + sofwareURL)
			fmt.Println(err)
			c.SendStatus(418)

		}
		doc, err := html.Parse(string(body))
		if err != nil {
			fmt.Println("FAILED TO GET HTML: " + sofwareURL)

			c.SendStatus(418)
		}

		//software-header
		if elementWithIDExists(doc, "software-header") {
			c.Send("SMH")
			//c.Send(software)
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
	return getElementById(n, id) == nil
}
