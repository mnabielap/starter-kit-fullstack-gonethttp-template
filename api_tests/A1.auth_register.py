import sys
import os
import time
sys.path.append(os.path.abspath(os.path.dirname(__file__)))
import utils

# Unique email to prevent 409 Conflict
unique_id = int(time.time())
email = f"user_{unique_id}@example.com"
password = "password123"

body = {
    "name": f"User {unique_id}",
    "email": email,
    "password": password
}

response = utils.send_and_print(
    url=f"{utils.BASE_URL}/auth/register",
    method="POST",
    body=body,
    output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json"
)

# Save credentials and tokens for next steps
if response.status_code == 201:
    data = response.json()
    
    root = data.get("data", data)
    
    if "tokens" in root:
        utils.save_config("access_token", root["tokens"]["access"]["token"])
        utils.save_config("refresh_token", root["tokens"]["refresh"]["token"])
        utils.save_config("current_user_email", email)
        utils.save_config("current_user_password", password)
        print(f"\n[INFO] Saved tokens for {email}")