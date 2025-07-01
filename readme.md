# readme

- SPA
- templ as a react replacement
- go backend
- htmx to extend html
- probably tailwind
- maybe typescript or jsdoc as a ts replacement
- maybe air

## setup

You can actually skip this step, as Go will automatically download dependencies when you run the code.
However, if you want to download them manually, you can run:

```bash
go mod download
```

## run

You can run the code with and get a development server on port 3000:

```bash
go tool templ generate && go run .
```

## build

I never tried it actually.

```bash
go tool templ generate && go build -o bin/app
```
