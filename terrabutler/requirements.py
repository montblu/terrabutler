from sys import exit

def check_requirements():
    """
    Check requirements before running the application.
    """
    from os import getenv, path
    from colorama import Fore

    if getenv("TERRABUTLER_ENABLE") != "true":
        print(Fore.YELLOW + "Terrabutler is not currently enabled on this"
              " folder. Please set 'TERRABUTLER_ENABLE' in your environment"
              " to true to enable it." + Fore.RESET)
        exit(1)
    root = getenv("TERRABUTLER_ROOT")
    if root is None or not path.exists(root):
        print(Fore.RED + "Terrabutler can't determine the root folder of"
              " your project or it doesn't exist.\nPlease set"
              " 'TERRABUTLER_ROOT' in your environment pointing"
              " to the root folder of your project." + Fore.RESET)
        exit(1)
    if not path.exists(root + "/configs/settings.yml"):
        print(Fore.RED + "Terrabutler can't find you settings file\nPlease"
              " create a 'settings.yml' file inside the 'configs' folder."
              + Fore.RESET)
        exit(1)
