# template-go-vercel

## Compile html

The /public folder is being used to serve all the static files.

To simplify the deployment, we want to compile all the go templates to html ahead of time

```
go run compiler/compiler.go
```

TODO: Write a watcher for this go file

## Compile css

Currently we are using tailwind, that compiles only on the compiled index.html

```
yarn tailwindcss -i templates/main.css -o public/style.css --watch
tailwind public/index.html > public/style.css
```

When compiling to deploy use:

```
yarn tailwindcss -i templates/main.css -o public/style.css --watch
```
