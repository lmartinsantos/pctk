include("resources:scripts/common")

pirates = actor {
    name = "men of low moral fiber (pirates)",
    size = {w=60, h=64},
    talkcolor = magenta,
    usepos = {x=90, y=128},
    usedir = LEFT
}

melee = room {
    background = "resources:backgrounds/Melee",
    collision_map = "resources:collision_maps/Melee_cm",
    objects = {
        bucket = object {
            class = APPLICABLE,
            name = "bucket",
            sprites = "resources:sprites/objects",
            pos = {x=260, y=120},
            hotspot = {x=250, y=100, w=20, h=20},
            usedir = RIGHT,
            usepos = {x=240, y=120},
            states = {
                default = {
                    anim = { 
                        { row = 6, delay = 1000, seq = {5} }
                    }
                },
                pickup = {}
            }
        },
        clock = object {
            name = "clock",
            hotspot = {x=150, y=25, w=24, h=18},
            usedir = UP,
            usepos = {x=161, y=116}
        }
    }
}

function melee:enter()
    local pirate1_dialog_props = { pos = {x=60, y=20}, color = magenta }
    local pirate2_dialog_props = { pos = {x=60, y=50}, color = yellow }
    local skipintro = true

    pirates:show{
        pos={x=38, y=137},         
        dir=RIGHT,
    }

    guybrush:show{
        pos={x=340, y=140}, 
        dir=LEFT,
    }
    
    music1:play()
    cricket:play()
    guybrush:walkto({x=290, y=140}).wait()
    if not skipintro then
        userputoff()
        cursoroff()

        guybrush:say("Hello, I'm Guybrush Threepwood,\nmighty pirate!").wait()
        sayline("**Oh no! This guy again!**", pirate1_dialog_props)
        guybrush:walkto({x=120, y=140}).wait()
        guybrush:say("I think I've lost the keys to my boat.").wait()
        guybrush:say("Have you seen any keys?", {delay=2000}).wait()
        sayline("Eeerrrr... Nope!", pirate1_dialog_props)
        sleep(2000)
        
        music2:play()
        guybrush:walkto({x=120, y=120}).wait()
        guybrush:say("Where can I find the keys?", {delay=1000}).wait()
        guybrush:walkto({x=120, y=140}).wait()
        guybrush:say("Ooooook...").wait()
        sleep(2000)
        guybrush:stand({dir = RIGHT}).wait()
        sleep(2000)
        guybrush:say("Ok, I will try the Scumm bar.").wait()
        guybrush:stand({dir = LEFT}).wait()
        guybrush:say("Thank you guys!").wait()
        cricket:play()
        guybrush:walkto({x=360, y=140}).wait()
        
        sayline("Oh, Jesus! I though he would\ntell again that stupid\ntale about LeChuck!", pirate1_dialog_props).wait()
        sleep(5000)
        sayline("Who has the keys?", pirate2_dialog_props).wait()
        sleep(1000)
        sayline("Me!", pirate1_dialog_props)
    end

    guybrush:select()
    userputon()
    cursoron()
end

function melee.objects.bucket:lookat()
    if self:owner() == guybrush then
        guybrush:say("It's a empty bucket.\nBut it's ALL MINE!")
    else
        guybrush:say("It's a empty bucket.")
    end
end

function melee.objects.bucket:give(to)
    if to == pirates then
        guybrush:say("I'd rather not. I am afraid\nthey'd get attached to it.")
    else
        DEFAULT.give(self)
    end
end

function melee.objects.bucket:pickup()
    cursoroff()
    guybrush:say("I don't know how this could help\nme to find the keys, but...").wait()
    guybrush:toinventory(self)
    cursoron()
end

function melee.objects.bucket:use(on)
    print("use bucket with", on.id)
    if on == melee.objects.clock then
        guybrush:say("Time flies, but I don't think\nI can gather it in the bucket.")
    elseif on == pirates then
        melee.objects.bucket:give(pirates)
    else
        DEFAULT.use(self, on)
    end
end

function melee.objects.clock:lookat()
    guybrush:say("It's weird. I have the feeling\nthat the time is not passing.").wait()
end

function melee.objects.clock:turnon()
    guybrush:say("Do I look like a watchmaker?").wait()
end

function melee.objects.clock:turnoff()
    guybrush:say("Well, I guess I couldn't be more off").wait()
end

function pirates:lookat() 
    guybrush:say("They didn't move since I arrived\nin Monkey Island I.").wait()
    guybrush:say("I guess they are waiting for\nsomething...").wait()
end

function pirates:talkto()
    guybrush:say("Now they are busy.\nI will not disturb them.")
end