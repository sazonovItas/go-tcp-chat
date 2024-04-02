import { toast } from "vue3-toastify";

class notificationSystem {
  public notify(level: string, msg: string) {
    switch (level) {
      case "success":
        toast.success(msg, {
          theme: "dark",
          autoClose: 3000,
        });
        break;
      case "warning":
        toast.warning(msg, {
          theme: "dark",
          autoClose: 3000,
        });
        break;
      case "error":
        toast.error(msg, {
          theme: "dark",
          autoClose: 3000,
        });
        break;
      default:
        toast.info(msg, {
          theme: "dark",
          autoClose: 3000,
        });
    }
  }
}

export const NotifySystem = new notificationSystem();

class responseToast {
  public notify(statusCode: number, status: string): void {
    if (statusCode >= 100 && statusCode < 200) {
      toast.info(statusCode + "\n" + status, {
        theme: "dark",
        autoClose: 3000,
      });
    } else if (statusCode >= 200 && statusCode < 300) {
      toast.success(statusCode + "\n" + status, {
        theme: "dark",
        autoClose: 3000,
      });
    } else if (statusCode >= 400 && statusCode < 500) {
      toast.warning(statusCode + "\n" + status, {
        theme: "dark",
        autoClose: 3000,
      });
    } else if (statusCode >= 500) {
      toast.error(statusCode + "\n" + status, {
        theme: "dark",
        autoClose: 3000,
      });
    }
  }
}

export const ResponseToast = new responseToast();
