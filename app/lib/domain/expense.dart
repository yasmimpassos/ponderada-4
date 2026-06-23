class Expense {
  final String id;
  final String groupId;
  final String paidBy;
  final String paidByName;
  final double amount;
  final String description;
  final String expenseDate;

  Expense({
    required this.id,
    required this.groupId,
    required this.paidBy,
    this.paidByName = '',
    required this.amount,
    required this.description,
    required this.expenseDate,
  });

  factory Expense.fromJson(Map<String, dynamic> json) {
    return Expense(
      id: json['id'] as String,
      groupId: json['group_id'] as String,
      paidBy: json['paid_by'] as String,
      paidByName: json['paid_by_name'] as String? ?? '',
      amount: (json['amount'] as num).toDouble(),
      description: json['description'] as String? ?? '',
      expenseDate: json['expense_date'] as String,
    );
  }
}

class Balance {
  final String userId;
  final String name;
  final double balance;

  Balance({required this.userId, required this.name, required this.balance});

  factory Balance.fromJson(Map<String, dynamic> json) {
    return Balance(
      userId: json['user_id'] as String,
      name: json['name'] as String? ?? json['user_name'] as String? ?? '',
      balance: (json['balance'] as num).toDouble(),
    );
  }
}

class OCRResult {
  final String description;
  final double amount;
  final String date;

  OCRResult({
    required this.description,
    required this.amount,
    required this.date,
  });

  factory OCRResult.fromJson(Map<String, dynamic> json) {
    return OCRResult(
      description: json['description'] as String? ?? '',
      amount: (json['amount'] as num?)?.toDouble() ?? 0.0,
      date: json['date'] as String? ?? '',
    );
  }
}
