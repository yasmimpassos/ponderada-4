import 'notification_service_stub.dart'
    if (dart.library.html) 'notification_service_web.dart'
    if (dart.library.io) 'notification_service_mobile.dart' as impl;

class NotificationService {
  static Future<void> init() => impl.init();

  static Future<void> showExpenseAdded(String description, double amount) =>
      impl.showExpenseAdded(description, amount);
}
