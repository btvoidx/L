-- https://github.com/btvoidx/L

function task.build()
  description "Builds L cli to 'bin' directory"
  depends { task.tidy, task.test, task.lint }
  sources { "**/*.go", "go.mod" }

  for GOOS, ext in pairs({ ["windows"] = "win.exe", ["linux"] = "linux", ["darwin"] = "macos" }) do
    print("Building for " .. GOOS)
    os.setenv("GOOS", GOOS)
    os.execute(string.format("go build -o bin/L_%s ./cmd/L", ext))

    -- if os.os() ~= "windows" and GOOS ~= "windows" then
    --   os.execute("chmod +x ")
    -- end
  end

  print("Builds done, binaries are in 'bin' directory")
end

function task.tidy()
  -- description 'Runs "go mod tidy"'
  os.execute("go mod tidy")
end
