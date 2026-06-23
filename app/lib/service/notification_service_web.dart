import 'dart:html' as html;

Future<void> init() async {
  if (html.Notification.supported) {
    await html.Notification.requestPermission();
  }
}

Future<void> showExpenseAdded(String description, double amount) async {
  if (!html.Notification.supported) return;
  if (html.Notification.permission == 'granted') {
    html.Notification(
      'Despesa adicionada',
      body: 'R\$ ${amount.toStringAsFixed(2)} — $description',
    );
  }
}
