import 'dart:convert';
import 'package:http/http.dart' as http;
import '../domain/group.dart';
import 'auth_service.dart';
import 'config_service.dart';

class GroupService {
  final AuthService _auth = AuthService();

  Future<List<Group>> getGroups() async {
    final headers = await _auth.authHeaders();
    final response = await http.get(
      Uri.parse('${ConfigService.baseUrl}/groups'),
      headers: headers,
    );

    if (response.statusCode == 200) {
      final List<dynamic> data = jsonDecode(response.body);
      return data.map((json) => Group.fromJson(json)).toList();
    }

    throw Exception('Erro ao carregar grupos');
  }

  Future<Group> createGroup(String name, String description) async {
    final headers = await _auth.authHeaders();
    final response = await http.post(
      Uri.parse('${ConfigService.baseUrl}/groups'),
      headers: headers,
      body: jsonEncode({'name': name, 'description': description}),
    );

    if (response.statusCode == 201) {
      return Group.fromJson(jsonDecode(response.body));
    }

    final error = jsonDecode(response.body);
    throw Exception(error['error'] ?? 'Erro ao criar grupo');
  }

  Future<void> addMember(String groupId, String email) async {
    final headers = await _auth.authHeaders();
    final response = await http.post(
      Uri.parse('${ConfigService.baseUrl}/groups/$groupId/members'),
      headers: headers,
      body: jsonEncode({'email': email}),
    );

    if (response.statusCode != 200) {
      final error = jsonDecode(response.body);
      throw Exception(error['error'] ?? 'Erro ao adicionar membro');
    }
  }

  Future<void> joinGroup(String groupId) async {
    final headers = await _auth.authHeaders();
    final response = await http.post(
      Uri.parse('${ConfigService.baseUrl}/groups/$groupId/join'),
      headers: headers,
    );

    if (response.statusCode != 200) {
      final error = jsonDecode(response.body);
      throw Exception(error['error'] ?? 'Erro ao entrar no grupo');
    }
  }

  Future<List<GroupMember>> getMembers(String groupId) async {
    final headers = await _auth.authHeaders();
    final response = await http.get(
      Uri.parse('${ConfigService.baseUrl}/groups/$groupId/members'),
      headers: headers,
    );

    if (response.statusCode == 200) {
      final List<dynamic> data = jsonDecode(response.body);
      return data.map((json) => GroupMember.fromJson(json)).toList();
    }

    throw Exception('Erro ao carregar membros');
  }
}
