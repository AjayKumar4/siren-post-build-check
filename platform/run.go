package platform

func Run() {
	downloadDirectory, file := Download("siren-platform-demo-data", "11.1.7")
	extractedDirectory, _ := Unzip(file, downloadDirectory)

	ok1 := make(chan bool)
	ok2 := make(chan bool)

	go StartInvestigate(extractedDirectory, ok1)
	go StartFederate(extractedDirectory, ok2)

	<-ok1
	<-ok2

}
