import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))
import utils

email = utils.load_config("current_user_email") or "admin@example.com"

body = {
    "email": email
}

utils.send_and_print(
    url=f"{utils.BASE_URL}/auth/forgot-password",
    method="POST",
    body=body,
    output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json"
)