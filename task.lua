-- https://github.com/btvoidx/L

local cmdpath = "./cmd/L"

function task.build()
  description "Builds L cli"
  depends { task.test, task.lint }

  os.execute("go build -o bin " .. cmdpath)
  print("Done")
end

function task.default()
  description "A default task!"
  print("Hello world!")
end
