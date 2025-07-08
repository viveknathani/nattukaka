package main

import "player/services"

func main() {
	containerService := services.NewContainerService()

	result, err := containerService.BuildImage(
		"tym",
		"https://github.com/viveknathani/teachyourselfmath",
		"master",
		"HEAD",
	)
	if err != nil {
		panic(err)
	}

	println(result)
}
