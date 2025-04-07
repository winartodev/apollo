package repositories

const (
	InsertUserDBQuery = `
		INSERT INTO users 
		    (
				 uuid, 
				 email, 
				 phone_number,
				 username,
				 first_name,
				 last_name,
				 profile_picture,
				 password,
				 refresh_token,
				 is_email_verified,
				 is_phone_verified,
				 created_at,
				 updated_at
			) VALUES (
						$1,  -- uuid
						$2,  -- email
						$3,  -- phone_number
						$4,  -- username
						$5,  -- first_name
						$6,  -- last_name
						$7,  -- profile_picture
						$8,  -- password
						$9,  -- refresh_token
						$10, -- is_email_verified
						$11, -- is_phone_verified
						$12, -- created_at
						$13  -- updated_at 
					) 
			  RETURNING id;
	`

	GetUserPasswordByEmailDBQuery = `
		SELECT 
		    password 
		FROM users 
		WHERE 
		    email = $1;
	`

	GetRefreshTokenByIDDBQuery = `
		Select 
			refresh_token
		FROM users
		where 
		    id = $1;
	`

	GetUserQueryDB = `
		SELECT 
			id,
			uuid, 
			email, 
			phone_number,
			username,
			first_name,
			last_name,
			profile_picture,
			is_email_verified,
			is_phone_verified,
			last_login,
			created_at,
			updated_at
		FROM users
	`

	UpdateRefreshTokenByIDDBQuery = `
		UPDATE users 
		SET 
		    refresh_token = $1 
		WHERE 
		    id = $2;
	`

	IsRefreshTokenIsExistsDBQuery = `
		SELECT EXISTS ( SELECT 1 FROM users WHERE id = $1 AND (refresh_token <> '' or refresh_token IS NOT NULL)) AS refresh_token_exists
	`

	IsUserExistDBQuery = `
		SELECT
			EXISTS (SELECT 1 FROM users WHERE username = $1) AS username_is_exists,
			EXISTS (SELECT 1 FROM users WHERE email = $2) as email_is_exists,
			EXISTS (SELECT 1 FROM users WHERE phone_number = $3) as phone_number_is_exists
	`
)
