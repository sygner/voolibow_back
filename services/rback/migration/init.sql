INSERT INTO permissions VALUES 
('user_account.list_sessions', 'List User Sessions'),
('user_user.list_users', 'List Users'),
('user_user.change_user_status', 'Change User Status'),
('user_user.delete_user', 'Delete User'),
('user_device.register_device', 'Register Device'),
('user_device.delete_device', 'Delete Device'),
('user_device.get_user_devices', 'Get User Devices');

INSERT INTO roles VALUES 
('owner','Owner Role');
('user','User Role');

INSERT INTO role_permissions VALUES
('owner', 'user_account.list_sessions'),
('owner', 'user_user.list_users'),
('owner', 'user_user.change_user_status'),
('owner', 'user_user.delete_user'),
('user', 'user_device.register_device'),
('user', 'user_device.delete_device'),
('user', 'user_device.get_user_devices');
