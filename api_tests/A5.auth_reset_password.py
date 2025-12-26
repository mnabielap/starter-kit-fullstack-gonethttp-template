import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))
import utils

# In a real scenario, you'd paste the token from the console log (since we mocked email in Go)
# For now, we simulate a call. It will likely fail without a valid token.
token = "PASTE_VALID_TOKEN_HERE_FROM_CONSOLE_LOG"
new_password = "newpassword123"

body = {
    "password": new_password
}

utils.send_and_print(
    url=f"{utils.BASE_URL}/auth/reset-password?token={token}",
    method="POST",
    body=body,
    output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json"
)