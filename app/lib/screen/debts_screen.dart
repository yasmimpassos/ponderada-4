import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import '../domain/group_balance.dart';
import '../service/expense_service.dart';

class DebtsScreen extends StatefulWidget {
  const DebtsScreen({super.key});

  @override
  State<DebtsScreen> createState() => _DebtsScreenState();
}

class _DebtsScreenState extends State<DebtsScreen> {
  final ExpenseService _expenseService = ExpenseService();
  List<PersonBalance> _balances = [];
  bool _loading = true;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    setState(() => _loading = true);
    try {
      final balances = await _expenseService.getPersonalBalances();
      if (mounted) setState(() => _balances = balances);
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(e.toString().replaceFirst('Exception: ', ''))),
        );
      }
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  void _showSettleDialog(PersonBalance balance) {
    final currencyFormat = NumberFormat.currency(locale: 'pt_BR', symbol: 'R\$');
    final amountController = TextEditingController(
      text: balance.balance.abs().toStringAsFixed(2),
    );

    final iOwe = balance.balance < 0;
    final title = iOwe
        ? 'Pagar ${balance.userName}'
        : 'Marcar recebido de ${balance.userName}';

    showDialog(
      context: context,
      builder: (context) {
        return AlertDialog(
          title: Text(title),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                iOwe
                    ? 'Você deve ${currencyFormat.format(balance.balance.abs())} para ${balance.userName}.'
                    : '${balance.userName} te deve ${currencyFormat.format(balance.balance)}.',
                style: const TextStyle(fontSize: 14),
              ),
              const SizedBox(height: 16),
              TextField(
                controller: amountController,
                keyboardType: const TextInputType.numberWithOptions(decimal: true),
                decoration: const InputDecoration(
                  labelText: 'Valor pago (R\$)',
                  border: OutlineInputBorder(),
                ),
              ),
            ],
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: const Text('Cancelar'),
            ),
            FilledButton(
              onPressed: () async {
                final amount = double.tryParse(
                  amountController.text.replaceAll(',', '.'),
                );
                if (amount == null || amount <= 0) return;
                Navigator.pop(context);

                try {
                  if (iOwe) {
                    await _expenseService.settle(balance.userId, amount);
                  } else {
                    await _expenseService.settleReceived(balance.userId, amount);
                  }
                  _load();
                } catch (e) {
                  if (!mounted) return;
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(content: Text(e.toString().replaceFirst('Exception: ', ''))),
                  );
                }
              },
              child: const Text('Confirmar'),
            ),
          ],
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    final currencyFormat = NumberFormat.currency(locale: 'pt_BR', symbol: 'R\$');

    final iOwe = _balances.where((b) => b.balance < 0).toList();
    final theyOwe = _balances.where((b) => b.balance > 0).toList();

    final totalOwed = iOwe.fold(0.0, (s, b) => s + b.balance.abs());
    final totalToReceive = theyOwe.fold(0.0, (s, b) => s + b.balance);

    return Scaffold(
      appBar: AppBar(title: const Text('Dívidas')),
      body: _loading
          ? const Center(child: CircularProgressIndicator())
          : RefreshIndicator(
              onRefresh: _load,
              child: _balances.isEmpty
                  ? const Center(
                      child: Text(
                        'Tudo quitado! Sem dívidas pendentes.',
                        style: TextStyle(fontSize: 16, color: Colors.grey),
                      ),
                    )
                  : ListView(
                      padding: const EdgeInsets.all(16),
                      children: [
                        _buildSummary(currencyFormat, totalOwed, totalToReceive),
                        const SizedBox(height: 20),
                        if (iOwe.isNotEmpty) ...[
                          _sectionHeader('Você deve', Colors.red),
                          const SizedBox(height: 8),
                          ...iOwe.map((b) => _PersonCard(
                                balance: b,
                                currencyFormat: currencyFormat,
                                onSettle: () => _showSettleDialog(b),
                              )),
                          const SizedBox(height: 20),
                        ],
                        if (theyOwe.isNotEmpty) ...[
                          _sectionHeader('Te devem', Colors.green),
                          const SizedBox(height: 8),
                          ...theyOwe.map((b) => _PersonCard(
                                balance: b,
                                currencyFormat: currencyFormat,
                                onSettle: () => _showSettleDialog(b),
                              )),
                        ],
                      ],
                    ),
            ),
    );
  }

  Widget _buildSummary(NumberFormat fmt, double owed, double toReceive) {
    final net = toReceive - owed;
    return Card(
      color: net >= 0 ? Colors.green.shade50 : Colors.red.shade50,
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            Text('Saldo geral',
                style: TextStyle(fontSize: 13, color: Colors.grey.shade600)),
            const SizedBox(height: 6),
            Text(
              net >= 0 ? '+${fmt.format(net)}' : fmt.format(net),
              style: TextStyle(
                fontSize: 26,
                fontWeight: FontWeight.bold,
                color: net >= 0 ? Colors.green.shade700 : Colors.red.shade700,
              ),
            ),
            const SizedBox(height: 12),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                _summaryItem('Você deve', fmt.format(owed), Colors.red),
                _summaryItem('Te devem', fmt.format(toReceive), Colors.green),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _summaryItem(String label, String value, Color color) {
    return Column(
      children: [
        Text(label, style: TextStyle(fontSize: 12, color: color)),
        Text(value,
            style: TextStyle(fontWeight: FontWeight.bold, color: color)),
      ],
    );
  }

  Widget _sectionHeader(String text, Color color) {
    return Text(text,
        style: TextStyle(
            fontSize: 16, fontWeight: FontWeight.bold, color: color));
  }
}

class _PersonCard extends StatelessWidget {
  final PersonBalance balance;
  final NumberFormat currencyFormat;
  final VoidCallback onSettle;

  const _PersonCard({
    required this.balance,
    required this.currencyFormat,
    required this.onSettle,
  });

  @override
  Widget build(BuildContext context) {
    final iOwe = balance.balance < 0;
    final color = iOwe ? Colors.red : Colors.green;
    final amount = balance.balance.abs();

    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: ListTile(
        leading: CircleAvatar(
          backgroundColor: color.withOpacity(0.15),
          child: Text(
            balance.userName[0].toUpperCase(),
            style: TextStyle(color: color, fontWeight: FontWeight.bold),
          ),
        ),
        title: Text(balance.userName,
            style: const TextStyle(fontWeight: FontWeight.w600)),
        subtitle: Text(
          iOwe
              ? 'Você deve ${currencyFormat.format(amount)}'
              : 'Te deve ${currencyFormat.format(amount)}',
          style: TextStyle(color: color),
        ),
        trailing: FilledButton.tonal(
          onPressed: onSettle,
          style: FilledButton.styleFrom(
            backgroundColor: color.withOpacity(0.1),
            foregroundColor: color,
          ),
          child: Text(iOwe ? 'Pagar' : 'Recebido'),
        ),
      ),
    );
  }
}
