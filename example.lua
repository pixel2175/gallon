gallon.before_request(function(req)
    -- Runs before every route.
    return req
end)

gallon.route("/", function(req)
    return "<h1>Hello</h1>", 200, {["Content-Type"] = "text/html"}
end)

gallon.route("GET /task/{task_num}", function(req)
    return "Task " .. req.path_value("task_num")
end)

gallon.route("/heavy", function(req)
    local n = 0
    for i = 1, 5000000 do n = n + i end
    return "done " .. n
end)

gallon.route("/redirect", function(req)
    return "<h1>Redirecting</h1>", 301
end)

gallon.error_handler(404, function(req)
    return "Not Found"
end)

gallon.error_handler(500, function(req)
    return "Internal Server Error"
end)

gallon.after_request(function(req)
    -- Runs after every route.
    return req
end)

gallon.jinja.set_template_dir("templates") -- default is templates/ too

-- Render templates/index.html.
gallon.route("/jinja/file", function(req)
    return gallon.jinja.render_template("index.html", {
        title = "Hello from Gonja",
        message = "Rendered from a template file"
    })
	, 200, {["Content-Type"] = "text/html"}
end)

-- Render a template directly from a string.
gallon.route("/jinja/string", function(req)
    return gallon.jinja.render_template_from_string(
        "<h1>{{ title }}</h1><p>{{ message }}</p>",
        {
            title = "Hello from Gonja",
            message = "Rendered from a string"
        }
    ), 200, {["Content-Type"] = "text/html"}
end)

gallon.run({host = "localhost", port = 8081, workers = 4, gc = 400})
