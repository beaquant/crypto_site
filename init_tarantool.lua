box.cfg{listen="127.0.0.1:3309", work_dir='./sessions', log_level=5, log='./tarantool.log'}
log = require('log')

memtx = false

if box.space.sessions == nil then
    if memtx then
        box.schema.space.create('sessions')
    else
        box.schema.space.create('sessions', {engine='vinyl'})
    end    
    box.space.sessions:create_index('primary', {parts={1, 'string'}})
    box.schema.user.grant('guest', 'read,write', 'space', 'sessions')
end