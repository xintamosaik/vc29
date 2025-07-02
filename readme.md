# readme

## plans

- [x] SPA
- [x] templ as a react replacement
- [x] go backend
- [x] htmx to extend html
- [ ] probably tailwind or maybe not?
- [ ] maybe typescript or jsdoc as a ts replacement
- [ ] maybe air for hot reloading

## workarounds (unplanned features)

- [x] esbuild to bundle js and css

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
