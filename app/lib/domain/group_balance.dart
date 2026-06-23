class PersonBalance {
  final String userId;
  final String userName;
  final double balance;

  PersonBalance({
    required this.userId,
    required this.userName,
    required this.balance,
  });

  factory PersonBalance.fromJson(Map<String, dynamic> json) {
    return PersonBalance(
      userId: json['user_id'] as String,
      userName: json['user_name'] as String,
      balance: (json['balance'] as num).toDouble(),
    );
  }
}
