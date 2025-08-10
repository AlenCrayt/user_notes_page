deploy:
	@mkdir -p build
	@cd code && go build
	@mv -f ./code/notes build
	@cp -r web_page build
