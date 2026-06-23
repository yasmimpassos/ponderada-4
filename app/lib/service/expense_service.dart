import 'dart:convert';
import 'dart:typed_data';
import 'package:http/http.dart' as http;
import '../domain/expense.dart';
import '../domain/group_balance.dart' show PersonBalance;
import 'auth_service.dart';
import 'config_service.dart';

class ExpenseService {
  final AuthService _auth = AuthService();

  Future<List<PersonBalance>> getPersonalBalances() async {
    final headers = await _auth.authHeaders();
    final response = await http.get(
      Uri.parse('${ConfigService.baseUrl}/balances'),
      headers: headers,
    );

    if (response.statusCode == 200) {
      final List<dynamic> data = jsonDecode(response.body);
      return data.map((json) => PersonBalance.fromJson(json)).toList();
    }

    throw Exception('Erro ao carregar dívidas');
  }

  Future<void> settle(String payeeId, double amount) async {
    final headers = await _auth.authHeaders();
    final response = await http.post(
      Uri.parse('${ConfigService.baseUrl}/settlements'),
      headers: headers,
      body: jsonEncode({'payee_id': payeeId, 'amount': amount}),
    );

    if (response.statusCode != 201) {
      final error = jsonDecode(response.body);
      throw Exception(error['error'] ?? 'Erro ao registrar pagamento');
    }
  }

  Future<void> settleReceived(String payerId, double amount) async {
    final headers = await _auth.authHeaders();
    final response = await http.post(
      Uri.parse('${ConfigService.baseUrl}/settlements/received'),
      headers: headers,
      body: jsonEncode({'payer_id': payerId, 'amount': amount}),
    );

    if (response.statusCode != 201) {
      final error = jsonDecode(response.body);
      throw Exception(error['error'] ?? 'Erro ao registrar recebimento');
    }
  }

  Future<List<Expense>> getGroupExpenses(String groupId) async {
    final headers = await _auth.authHeaders();
    final response = await http.get(
      Uri.parse('${ConfigService.baseUrl}/groups/$groupId/expenses'),
      headers: headers,
    );

    if (response.statusCode == 200) {
      final List<dynamic> data = jsonDecode(response.body);
      return data.map((json) => Expense.fromJson(json)).toList();
    }

    throw Exception('Erro ao carregar despesas');
  }

  Future<List<Balance>> getGroupBalances(String groupId) async {
    final headers = await _auth.authHeaders();
    final response = await http.get(
      Uri.parse('${ConfigService.baseUrl}/groups/$groupId/balances'),
      headers: headers,
    );

    if (response.statusCode == 200) {
      final List<dynamic> data = jsonDecode(response.body);
      return data.map((json) => Balance.fromJson(json)).toList();
    }

    throw Exception('Erro ao carregar balanços');
  }

  Future<Expense> createExpense({
    required String groupId,
    required String description,
    required double amount,
    required String expenseDate,
    required List<String> splitUserIds,
  }) async {
    final headers = await _auth.authHeaders();
    final response = await http.post(
      Uri.parse('${ConfigService.baseUrl}/groups/$groupId/expenses'),
      headers: headers,
      body: jsonEncode({
        'description': description,
        'amount': amount,
        'expense_date': expenseDate,
        'split_user_ids': splitUserIds,
      }),
    );

    if (response.statusCode == 201) {
      return Expense.fromJson(jsonDecode(response.body));
    }

    final error = jsonDecode(response.body);
    throw Exception(error['error'] ?? 'Erro ao criar despesa');
  }

  Future<void> deleteExpense(String groupId, String expenseId) async {
    final headers = await _auth.authHeaders();
    final response = await http.delete(
      Uri.parse('${ConfigService.baseUrl}/groups/$groupId/expenses/$expenseId'),
      headers: headers,
    );

    if (response.statusCode != 200) {
      final error = jsonDecode(response.body);
      throw Exception(error['error'] ?? 'Erro ao deletar despesa');
    }
  }

  Future<OCRResult> processOCR(Uint8List imageBytes, String filename) async {
    final headers = await _auth.authHeaders();
    headers.remove('Content-Type');

    final request = http.MultipartRequest(
      'POST',
      Uri.parse('${ConfigService.baseUrl}/ocr'),
    );
    request.headers.addAll(headers);
    request.files.add(http.MultipartFile.fromBytes(
      'image',
      imageBytes,
      filename: filename,
    ));

    final streamedResponse = await request.send();
    final response = await http.Response.fromStream(streamedResponse);

    if (response.statusCode == 200) {
      return OCRResult.fromJson(jsonDecode(response.body));
    }

    throw Exception('Erro ao processar imagem');
  }
}
