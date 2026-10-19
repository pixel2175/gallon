gallon.before_request(function(req)
	return req
    -- runs before every route
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
    -- runs after every route
	return req
end)

gallon.run({host = "localhost", port = 8081, workers = 4, gc=200})

