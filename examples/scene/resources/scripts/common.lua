guybrush = actor { }

music1 = music { ref = "resources:audio/OnTheHill" }
music2 = music { ref = "resources:audio/GuitarNoodling" }
cricket = sound { ref = "resources:audio/Cricket" }

magenta = { r = 0xAA, g = 0x00, b = 0xAA }
yellow = { r = 0xFF, g = 0xFF, b = 0x55 }

function default.pickup()
    guybrush:say("I can't pick that up.")
end

function default.use()
    guybrush:say("I can't use that.")
end

function default.open()
    guybrush:say("I can't open that.")
end

function default.close()
    guybrush:say("I can't close that.")
end

function default.pull()
    guybrush:say("I can't pull that.")
end

function default.push()
    guybrush:say("I can't push that.")
end

function default.talkto()
    guybrush:say("I can't talk to that.")
end

function default.lookat()
    guybrush:say("There is nothing special about that.")
end

function default.turnon()
    guybrush:say("I can't turn that on.")
end

function default.turnoff()
    guybrush:say("I can't turn that off.")
end

function default.give()
    guybrush:say("I can't give that.")
end

function default.walkto()    
end
