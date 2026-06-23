class Group {
  final String id;
  final String name;
  final String description;
  final String createdBy;

  Group({
    required this.id,
    required this.name,
    this.description = '',
    required this.createdBy,
  });

  factory Group.fromJson(Map<String, dynamic> json) {
    return Group(
      id: json['id'] as String,
      name: json['name'] as String,
      description: json['description'] as String? ?? '',
      createdBy: json['created_by'] as String,
    );
  }
}

class GroupMember {
  final String userId;
  final String name;
  final String email;

  GroupMember({
    required this.userId,
    required this.name,
    required this.email,
  });

  factory GroupMember.fromJson(Map<String, dynamic> json) {
    return GroupMember(
      userId: json['user_id'] as String,
      name: json['name'] as String,
      email: json['email'] as String,
    );
  }
}
