def check_requirements():
    """
    Check requirements before running the application.
    """
    from os import getenv, path
    from colorama import Fore

    if getenv("TERRABUTLER_ENABLE") != "true":
        print(Fore.YELLOW + "Terrabutler is running outside of a project"
              " folder or Please set 'TERRABUTLER_ENABLE' in your environment."
              "\nExiting...")
        exit(1)
    root = getenv("TERRABUTLER_ROOT")
    if root is None or not path.exists(root):
        print(Fore.RED + "Terrabutler can't determine the root folder of"
              " your project or it doesn't exist.\nPlease set"
              " 'TERRABUTLER_ROOT' in your environment pointing"
              " to the root folder of your project.")
        exit(1)