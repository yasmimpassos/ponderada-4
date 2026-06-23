import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';
import 'config_service.dart';

class AuthService {
  static const String _tokenKey = 'token';

  Future<String?> getToken() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString(_tokenKey);
  }

  Future<void> saveToken(String token) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_tokenKey, token);
  }

  Future<void> clearToken() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove(_tokenKey);
  }

  Future<Map<String, String>> authHeaders() async {
    final token = await getToken();
    return {
      'Content-Type': 'application/json',
      if (token != null) 'Authorization': 'Bearer $token',
    };
  }

  Future<void> register(String name, String email, String password) async {
    final response = await http.post(
      Uri.parse('${ConfigService.baseUrl}/register'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({'name': name, 'email': email, 'password': password}),
    );

    if (response.statusCode == 201) {
      final data = jsonDecode(response.body);
      await saveToken(data['token'] as String);
      return;
    }

    final error = jsonDecode(response.body);
    throw Exception(error['error'] ?? 'Erro ao registrar');
  }

  Future<void> login(String email, String password) async {
    final response = await http.post(
      Uri.parse('${ConfigService.baseUrl}/login'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({'email': email, 'password': password}),
    );

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      await saveToken(data['token'] as String);
      return;
    }

    final error = jsonDecode(response.body);
    throw Exception(error['error'] ?? 'Email ou senha incorretos');
  }

  Future<void> logout() async {
    await clearToken();
  }
}
