import 'package:flutter_local_notifications/flutter_local_notifications.dart';

final _plugin = FlutterLocalNotificationsPlugin();

Future<void> init() async {
  const android = AndroidInitializationSettings('@mipmap/ic_launcher');
  const ios = DarwinInitializationSettings();
  await _plugin.initialize(
    const InitializationSettings(android: android, iOS: ios),
  );
}

Future<void> showExpenseAdded(String description, double amount) async {
  const details = NotificationDetails(
    android: AndroidNotificationDetails(
      'expenses',
      'Despesas',
      channelDescription: 'Notificações de despesas',
      importance: Importance.high,
      priority: Priority.high,
    ),
    iOS: DarwinNotificationDetails(),
  );
  await _plugin.show(
    0,
    'Despesa adicionada',
    'R\$ ${amount.toStringAsFixed(2)} — $description',
    details,
  );
}
