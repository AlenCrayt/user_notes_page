deploy:
	@mkdir -p build
	@cd code && go build
	@mv -f ./code/notes build
	@cp -r web_page build

tailwind:
	@npx @tailwindcss/cli -i ./web_page/input.css -o ./web_page/output.css --watch
