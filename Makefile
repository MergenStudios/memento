windows:
	go install memento.go

	go build service.go
	sc.exe create mementoPermSourcesUpdater start="auto" binPath="./service.exe"