box.cfg {
    -- wal_dir = nil;

    -- An absolute path to directory where snapshot (.snap) files are stored.
    -- If not specified, defaults to /var/lib/tarantool/INSTANCE
    -- snap_dir = nil;

    -- An absolute path to directory where vinyl files are stored.
    -- If not specified, defaults to /var/lib/tarantool/INSTANCE
    -- vinyl_dir = nil;

    -- The read/write data port number or URI
    -- Has no default value, so must be specified if
    -- connections will occur from remote clients
    -- that do not use “admin address”
    -- listen = 'localhost:3301';
    listen = '*:3301';

    -- Inject the given string into server process title
    -- custom_proc_title = 'example';

    -------------------------
    -- Storage configuration
    -------------------------

    -- How much memory Tarantool allocates
    -- to actually store tuples, in gigabytes

    slab_alloc_arena = 0.5;

    -- Size of the smallest allocation unit
    -- It can be tuned up if most of the tuples are not so small
    slab_alloc_minimal = 16;

    -- Size of the largest allocation unit
    -- It can be tuned up if it is necessary to store large tuples
    slab_alloc_maximal = 1048576;

    -- Use slab_alloc_factor as the multiplier for computing
    -- the sizes of memory chunks that tuples are stored in
    slab_alloc_factor = 1.06;

    -------------------
    -- Snapshot daemon
    -------------------

    -- The interval between actions by the snapshot daemon, in seconds
    snapshot_period = 5;

    -- The maximum number of snapshots that the snapshot daemon maintans
    snapshot_count = 6;

    --------------------------------

    -- Abort if there is an error while reading
    -- the snapshot file (at server start)
    panic_on_snap_error = true;

    -- Abort if there is an error while reading a write-ahead
    -- log file (at server start or to relay to a replica)
    panic_on_wal_error = true;

    -- How many log records to store in a single write-ahead log file
    rows_per_wal = 5000000;

    -- Reduce the throttling effect of box.snapshot() on
    -- INSERT/UPDATE/DELETE performance by setting a limit
    -- on how many megabytes per second it can write to disk
    snap_io_rate_limit = nil;

    -- Specify fiber-WAL-disk synchronization mode as:
    -- "none": write-ahead log is not maintained;
    -- "write": fibers wait for their data to be written to the write-ahead log;
    -- "fsync": fibers wait for their data, fsync follows each write;
    wal_mode = "none";

    -- Number of seconds between periodic scans of the write-ahead-log

    wal_dir_rescan_delay = 2.0;

    ---------------
    -- Replication
    ---------------

    -- The server is considered to be a Tarantool replica
    -- it will try to connect to the master
    -- which replication_source specifies with a URI
    -- for example konstantin:secret_password@tarantool.org:3301
    -- by default username is "guest"
    -- replication_source="127.0.0.1:3102";

    --------------
    -- Networking
    --------------

    -- The server will sleep for io_collect_interval seconds
    -- between iterations of the event loop
    io_collect_interval = nil;

    -- The size of the read-ahead buffer associated with a client connection
    readahead = 16320;

    ----------
    -- Logging
    ----------

    -- How verbose the logging is. There are six log verbosity classes:
    -- 1 – SYSERROR
    -- 2 – ERROR
    -- 3 – CRITICAL
    -- 4 – WARNING
    -- 5 – INFO
    -- 6 – DEBUG
    log_level = 5;

    -- By default, the log is sent to /var/log/tarantool/INSTANCE.log
    -- If logger is specified, the log is sent to the file named in the string
    -- logger = "example.log";

    -- If true, tarantool does not block on the log file descriptor
    -- when it’s not ready for write, and drops the message instead
    logger_nonblock = true;

    -- If processing a request takes longer than
    -- the given value (in seconds), warn about it in the log
    too_long_threshold = 0.5;
}

local function bootstrap()
    local space = box.schema.create_space('sessions')
    space:create_index('primary', {type = 'hash', parts = { 1, 'string'}})
    -- name, email, password, audios, bears
    local profilespace = box.schema.create_space('profile')
    profilespace:create_index('primary', {type = 'hash', parts = {1, 'string'}})
    --
    local toyspace = box.schema.create_space('toy')
    toyspace:create_index('primary', {type = 'hash', parts = {1, 'unsigned'}})
    --id, category, name, price, duration, descriprion, author
    local audiospace = box.schema.create_space('audio')
    audiospace:create_index('primary', {type = 'tree', parts = {1, 'unsigned'}})
    audiospace:create_index('cat_name', {type='tree',unique=true, parts = {2,'unsigned',3,'string'}})
    audiospace:create_index('cat_price_name', {type='tree',unique=true, parts = {2,'unsigned',4,'unsigned',3,'string'}})
    audiospace:create_index('cat_duration_name', {type='tree',unique=true, parts = {2,'unsigned',5,'string',3,'string'}})
    audiospace:create_index('name', {type='tree',unique=true, parts = {3,'string'}})
    audiospace:create_index('price_name', {type='tree',unique=true, parts = {4,'unsigned',3,'string'}})
    audiospace:create_index('duration_name', {type='tree',unique=true, parts = {5,'string',3,'string'}})
    box.schema.func.create('getProfile')
    box.schema.func.create('createProfile')
    box.schema.func.create('isLogined')
    box.schema.func.create('getAllStories')
    box.schema.func.create('getAllCategoryStories')
    box.schema.func.create('getUserStories')
    box.schema.func.create('buyStory')
    box.schema.func.create('findStory')
    -- Comment this if you need fine grained access control (without it, guest
    -- will have access to everything)
    -- box.schema.user.grant('goClient', 'read,write,execute', 'universe')

    -- Keep things safe by default
    box.schema.user.create('goClient', { password = 'TeddyTarantoolS1cret' })
    -- box.schema.user.grant('giClient', 'replication')
    box.schema.user.grant('goClient', 'read,write,execute', 'space', 'sessions')
    box.schema.user.grant('goClient', 'read,write,execute', 'space', 'profile')
    box.schema.user.grant('goClient', 'read,write,execute', 'space', 'toy')
    box.schema.user.grant('goClient', 'read,write,execute', 'space', 'audio')
    box.schema.user.grant('goClient', 'execute', 'function', 'getProfile')
    box.schema.user.grant('goClient', 'execute', 'function', 'createProfile')
    box.schema.user.grant('goClient', 'execute', 'function', 'isLogined')
    box.schema.user.grant('goClient', 'execute', 'function', 'getAllStories')
    box.schema.user.grant('goClient', 'execute', 'function', 'getAllCategoryStories')
    box.schema.user.grant('goClient', 'execute', 'function', 'getUserStories')
    box.schema.user.grant('goClient', 'execute', 'function', 'buyStory')
    box.schema.user.grant('goClient', 'execute', 'function', 'findStory')
end

-- for first run create a space and add set up grants
box.once('sessions-1.1', bootstrap)


function getProfile(sid)
    session = box.space.sessions:select { sid }
    if #session ~= 0 then
        name = session[1][2]['name']
        profile = box.space.profile:select { name }
        return profile
    end
    return error("no such session")
end

function createProfile(name, email, password)
    t = box.space.profile:insert { name, email, password }
    return t
end

function isLogined(sid)
    session = box.space.sessions:select{sid}
    if #session ~= 0 then
        name = session[1][2]['name']
        isExists = box.space.profile:count{name}
        if isExists == 1 then
            return true
        end
        return error("no such user")
    end
end

function getAllStories(offset, limit, order, ordertype)
    if ordertype == 1 then
        localiterat = box.index.GE
    else
        localiterat = box.index.LE
    end
    if order == 0 then
        stories = box.space.audio.index.name:select({},{iterator=localiterat, offset=offset, limit=limit})
    elseif order == 1 then
        stories = box.space.audio.index.price_name:select({},{iterator=localiterat, offset=offset, limit=limit})
    elseif order == 2 then
        stories = box.space.audio.index.duration_name:select({},{iterator=localiterat, offset=offset, limit=limit})
    end
    return stories
end

function getAllCategoryStories(category, offset, limit, order, ordertype)
    if ordertype == 1 then
        localiterat = box.index.EQ
    else
        localiterat = box.index.REQ
    end
    if order == 0 then
        stories = box.space.audio.index.cat_name:select({category},{iterator=localiterat, offset=offset, limit=limit})
    elseif order == 1 then
        stories = box.space.audio.index.cat_price_name:select({category},{iterator=localiterat, offset=offset, limit=limit})
    elseif order == 2 then
        stories = box.space.audio.index.cat_duration_name:select({category},{iterator=localiterat, offset=offset, limit=limit})
    end
    return stories
end

function getUserStories(userLogin)
    user = box.space.profile:get{userLogin}
    stories = {}
    for i, v in ipairs(user[4]) do
        tempstory = box.space.audio:get{v}
        table.insert(stories,tempstory)
    end
    return stories
end

function has_value (tab, val)
    for index, value in ipairs (tab) do
        if value == val then
            return true
        end
    end
    return false
end

function buyStory(userLogin, storyid)
    user = box.space.profile:get{userLogin}
    story = box.space.audio:get{storyid}
    if user == nil then
        return error("no such user")
    end
    if story == nil then
        return error("no such story")
    end
    stories = user[4]
    if has_value(stories, storyid) then
        return error("already bought")
    end
    table.insert(stories, storyid)
    table.sort(stories)
    box.space.profile:update(userLogin, {{'=', 4, stories}})
    return stories
end

-- Return the minimum of three elements
function min(a, b, c)
    return math.min(math.min(a, b), c)
end

-- Creates a 2D matrix
function matrix(row,col)
    local m = {}
    for i = 1,row do m[i] = {}
    for j = 1,col do m[i][j] = 0 end
    end
    return m
end

-- Calculates the Levenshtein distance between two strings
function lev_iter_based(strA,strB)
    local M = matrix(#strA+1,#strB+1)
    local i,j,cost
    local row,col = #M,#M[1]
    for i = 1,row do M[i][1] = i-1 end
    for j = 1,col do M[1][j] = j-1 end
    for i = 2,row do
        for j = 2,col do
            if (strA:sub(i-1,i-1) == strB:sub(j-1,j-1)) then cost = 0
            else cost = 1
            end
            M[i][j] = min(M[i-1][j]+1,M[i][j-1]+1,M[i-1][j-1]+cost)
        end
    end
    return M[row][col]
end

function findStory(str)
    stories = {}
    for k, v in s:pairs() do
        lev_dist = lev_iter_based(v[3]:lower(), str:lower())
        if lev_dist < 3 then
            table.insert(stories,v)
        end
    end
    return stories
end

function addStory(name, description, author, duration, price, category)
    local s = box.space.audio:auto_increment{category, name, price,duration, description, author}
    return s
end

function addSubStory(id)
    local s = box.space.audio:get{id}
    if s == null then
        return error("no such story")
    end
    local arr = s[17]

    return s
end

function getSubStories(id)
    return json.encode(box.space.audio:get(id)[17])
end