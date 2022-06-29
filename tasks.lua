-- https://github.com/btvoidx/L

function task.build()
  description "Builds L cli to 'bin' directory"
  depends { task.tidy, task.test, task.lint }
  sources { "**/*.go", "go.mod" }

  os.execute("go build -o bin ./cmd/L")
  print("Build done, binaries are in 'bin' directory")
end

function task.tidy()
  description 'Runs "go mod tidy"'
  os.execute("go mod tidy")
end
