-- https://github.com/btvoidx/L

function task.build()
  description "Builds L cli"
  depends { task.tidy, task.test, task.lint }
  sources { "*.go", "go.mod" }

  os.execute("go build -o bin ./cmd/L")
  print("Build done, binaries are in 'bin' folder")
end

function task.tidy()
  os.execute("go mod tidy")
end

function task.default()
  description "A default task!"
  print("Hello world!")
end
