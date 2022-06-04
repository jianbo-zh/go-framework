

swagger.%:
	swag init -d ./cmd/${*} -o ./cmd/${*}/doc/

trans.%:
	@i18n -extract.dir=./cmd/${*} -output.dir=./cmd/${*}/translations extract
	@i18n -output.dir=./cmd/${*}/translations generator