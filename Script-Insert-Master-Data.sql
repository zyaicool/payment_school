--MASTER DATA ROLE
--list role
-- 1. super admin
-- 2. orang tua/siswa
-- 3. tata usaha
-- 4. kasir

-- query insert data role
INSERT INTO public.roles
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, "name")
VALUES(now(), 0, now(), 0, null, null, 'Super Admin');

INSERT INTO public.roles
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, "name")
VALUES(now(), 0, now(), 0, null, null, 'Orang Tua/Siswa');

INSERT INTO public.roles
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, "name")
VALUES(now(), 0, now(), 0, null, null, 'Tata Usaha');

INSERT INTO public.roles
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, "name")
VALUES(now(), 0, now(), 0, null, null, 'Kasir');

INSERT INTO public.roles
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, "name")
VALUES(now(), 0, now(), 0, null, null, 'Admin');

--MASTER DATA ROLE MATRIX
--list akses matrix role
-- 1. super admin
--	1.1. login					: R
--	1.2. registrasi				: CR
--	1.3. daftar akun			: CRUD
--	1.4. tambah akun			: CR
--	1.5. edit akun				: UR
--	1.6. dashboard siswa		: CRUD
--	1.7. tambah siswa baru		: CR
--	1.8. edit informasi siswa	: UR
--
-- 2. orang tua/siswa
--	2.1. login					: R
--	2.2. registrasi				: C
--	2.3. change password		: U
--	2.4. dashboard siswa		: CRU
--	2.5. edit informasi siswa	: U
--	2.6. daftar tagihan			: 
--
-- 3. tata usaha
--	3.1. login					: R
--	3.2. dashboard jenis tagihan:
--	3.3. tambah jenis tagihan	:
--	3.4. dashboard kelas		:
--
-- 4. kasir
--	4.1. login					: R
--	4.2. registrasi				:
--
-- keretangan:
-- 1. C = Create
-- 2. R = Read
-- 3. U = Update
-- 4. D = Delete

-- query insert data role matrix
-- role super admin
INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 1, 'Login', 'A002', false, true, false, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 1, 'Registrasi', 'A001', true, false, false, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 1, 'User Management', 'A003', true, true, true, true);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 1, 'Tambah Akun User', 'A004', true, false, false, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 1, 'Edit Akun User', 'A005', false, false, true, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 1, 'Dashboard Master Siswa', 'A006', true, true, true, true);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 1, 'Tambah Master Siswa', 'A007', true, false, false, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 1, 'Edit Master Siswa', 'A008', false, false, true, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 1, 'Homepage', 'A020', false, true, false, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 1, 'Forget Password', 'A021', false, false, true, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 1, 'Lihat Informasi Siswa', 'A022', false, true, false, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 1, 'Dashboard Pembayaran Tagihan', 'A023', false, true, true, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 1, 'Detail Pembayaran Tagihan', 'A024', false, false, true, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 1, 'Dashboard Riwayat Tagihan', 'A025', false, true, false, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 1, 'Detail Riwayat Tagihan', 'A026', false, true, false, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 1, 'tambah siswa baru', 'A007', true, false, false, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 1, 'edit informasi siswa', 'A008', false, false, true, false);

-- role orang tua/siswa
INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 2, 'Login', 'A002', false, true, false, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 2, 'Homepage', 'A020', false, true, false, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 2, 'Forget Password', 'A021', false, false, true, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 2, 'Registrasi', 'A001', true, false, false, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 2, 'Dashboard Siswa', 'A006', true, true, true, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 2, 'Tambah Siswa Baru', 'A007', true, false, false, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 2, 'Edit Informasi Siswa', 'A008', false, false, true, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 2, 'Lihat Informasi Siswa', 'A022', false, true, false, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 2, 'Dashboard Pembayaran Tagihan', 'A023', false, true, true, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 2, 'Detail Pembayaran Tagihan', 'A024', false, false, true, false);

-- role tata usaha
INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 3, 'Login', 'A002', false, true, false, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 3, 'Homepage', 'A020', false, true, false, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 3, 'Forget Password', 'A021', false, false, true, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 3, 'Dashboard Master Tagihan', 'A010', true, true, true, true);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 3, 'Tambah Data Tagihan', 'A011', true, false, false, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 3, 'edit jenis tagihan', 'A012', false, false, true, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 3, 'Dashboard Riwayat Pembayaran', 'A025', false, true, false, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 3, 'Detail Riwayat Pembayaran', 'A026', false, true, false, false);

-- role kasir

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 4, 'Login', 'A002', false, true, false, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 4, 'Homepage', 'A020', false, true, false, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 4, 'Forget Password', 'A021', false, false, true, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 4, 'Dashboard Pembayaran Tagihan', 'A023', false, true, false, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 4, 'Detail Pembayaran Tagihan', 'A024', false, false, true, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 4, 'Dashboard Riwayat Pembayaran', 'A025', false, true, false, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 4, 'Detail Riwayat Pembayaran', 'A026', false, true, false, false);

-- role admin sekolah
INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 5, 'Login', 'A002', false, true, false, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 5, 'Homepage', 'A020', false, true, false, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 5, 'Forget Password', 'A021', false, false, true, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 5, 'User Management', 'A003', true, true, true, true);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 5, 'Tambah Akun User', 'A004', true, false, false, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 5, 'Edit Akun User', 'A005', false, false, true, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 5, 'Dashboard Master Siswa', 'A006', true, true, true, true);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 5, 'Tambah Master Siswa', 'A007', true, false, false, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 5, 'Edit Master Siswa', 'A008', false, false, true, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 5, 'Lihat Informasi Siswa', 'A022', false, true, false, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 5, 'Dashboard Pembayaran Tagihan', 'A023', false, true, true, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 5, 'Detail Pembayaran Tagihan', 'A024', false, false, true, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 5, 'Dashboard Riwayat Pembayaran', 'A025', false, true, false, false);

INSERT INTO public.role_matrices
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, page_name, page_code, is_create, is_read, is_update, is_delete)
VALUES(now(), 0, now(), 0, null, null, 5, 'Detail Riwayat Pembayaran', 'A026', false, true, false, false);


-- insert data user super admin
INSERT INTO public.users
(created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, role_id, email, "password", "username", is_verification, is_block)
VALUES(now(), 0, now(), 0, null, null, 1, 'admin@gmail.com', '$2a$10$DxqnRmx94sPvti7Dg6Nhr.AuRHELaiirMteqalQ0OPAbi9.UWwcby', 'admin1', true, false);