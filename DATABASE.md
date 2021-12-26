## DATABASE
---

To make the new DB accessible on the network.  Edit /some/path/postgresql.conf and make the following changes

    #------------------------------------------------------------------------------
    # CONNECTIONS AND AUTHENTICATION
    #------------------------------------------------------------------------------

    # - Connection Settings -

    listen_addresses = '*'                                                   <= ADD
    #listen_addresses = 'localhost'		# what IP address(es) to listen on;
					# comma-separated list of addresses;
					# defaults to 'localhost'; use '*' for all
					# (change requires restart)


Give the new admin and user access to the db.  Edit /some/path/pg_hba.conf

    # Database administrative login by Unix domain socket
    local   all             postgres                                peer 

    # TYPE  DATABASE        USER            ADDRESS                 METHOD

    # "local" is for Unix domain socket connections only
    local   mynewdb         all                                     md5      <= ADD LINE
    local   all             all                                     peer
    # IPv4 local connections:
    host    mynewdb        newdbuser         0.0.0.0/0              md5      <= ADD LINE
    host    all             all             127.0.0.1/32            md5     


When the edits are complete, restart PostgreSQL

    > systemctl restart postgresql