import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:intl/intl.dart';
import '../domain/group.dart';
import '../domain/expense.dart';
import '../service/expense_service.dart';
import '../service/group_service.dart';
import 'new_expense_screen.dart';

class GroupScreen extends StatefulWidget {
  final Group group;

  const GroupScreen({super.key, required this.group});

  @override
  State<GroupScreen> createState() => _GroupScreenState();
}

class _GroupScreenState extends State<GroupScreen>
    with SingleTickerProviderStateMixin {
  final ExpenseService _expenseService = ExpenseService();
  final GroupService _groupService = GroupService();
  late TabController _tabController;

  List<Expense> _expenses = [];
  List<Balance> _balances = [];
  List<GroupMember> _members = [];
  bool _loading = true;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
    _loadData();
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  Future<void> _loadData() async {
    setState(() => _loading = true);
    try {
      final results = await Future.wait([
        _expenseService.getGroupExpenses(widget.group.id),
        _expenseService.getGroupBalances(widget.group.id),
        _groupService.getMembers(widget.group.id),
      ]);
      if (mounted) {
        setState(() {
          _expenses = results[0] as List<Expense>;
          _balances = results[1] as List<Balance>;
          _members = results[2] as List<GroupMember>;
        });
      }
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

  Future<void> _deleteExpense(String expenseId) async {
    try {
      await _expenseService.deleteExpense(widget.group.id, expenseId);
      _loadData();
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(e.toString().replaceFirst('Exception: ', ''))),
      );
    }
  }

  void _shareInviteLink() {
    final uri = Uri.base;
    final inviteLink = '${uri.scheme}://${uri.host}${uri.port != 80 && uri.port != 443 ? ':${uri.port}' : ''}?join=${widget.group.id}';

    showDialog(
      context: context,
      builder: (context) {
        return AlertDialog(
          title: const Text('Link de convite'),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text('Compartilhe este link para convidar alguém:'),
              const SizedBox(height: 12),
              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: Colors.grey.shade100,
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Text(
                  inviteLink,
                  style: const TextStyle(fontSize: 13),
                ),
              ),
            ],
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: const Text('Fechar'),
            ),
            FilledButton.icon(
              onPressed: () {
                Clipboard.setData(ClipboardData(text: inviteLink));
                Navigator.pop(context);
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(content: Text('Link copiado!')),
                );
              },
              icon: const Icon(Icons.copy),
              label: const Text('Copiar'),
            ),
          ],
        );
      },
    );
  }

  void _showAddMemberDialog() {
    final emailController = TextEditingController();

    showDialog(
      context: context,
      builder: (context) {
        return AlertDialog(
          title: const Text('Adicionar membro'),
          content: TextField(
            controller: emailController,
            keyboardType: TextInputType.emailAddress,
            decoration: const InputDecoration(
              labelText: 'Email do usuário',
              border: OutlineInputBorder(),
            ),
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: const Text('Cancelar'),
            ),
            FilledButton(
              onPressed: () async {
                if (emailController.text.isEmpty) return;
                Navigator.pop(context);
                try {
                  await _groupService.addMember(
                    widget.group.id,
                    emailController.text.trim(),
                  );
                  _loadData();
                } catch (e) {
                  if (!mounted) return;
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(content: Text(e.toString().replaceFirst('Exception: ', ''))),
                  );
                }
              },
              child: const Text('Adicionar'),
            ),
          ],
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.group.name),
        bottom: TabBar(
          controller: _tabController,
          tabs: const [
            Tab(text: 'Despesas'),
            Tab(text: 'Balanço'),
            Tab(text: 'Membros'),
          ],
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.share),
            tooltip: 'Compartilhar link de convite',
            onPressed: _shareInviteLink,
          ),
          IconButton(
            icon: const Icon(Icons.person_add),
            onPressed: _showAddMemberDialog,
          ),
        ],
      ),
      body: _loading
          ? const Center(child: CircularProgressIndicator())
          : TabBarView(
              controller: _tabController,
              children: [
                _buildExpensesTab(),
                _buildBalancesTab(),
                _buildMembersTab(),
              ],
            ),
      floatingActionButton: FloatingActionButton(
        onPressed: () async {
          await Navigator.of(context).push(
            MaterialPageRoute(
              builder: (_) => NewExpenseScreen(
                groupId: widget.group.id,
                members: _members,
              ),
            ),
          );
          _loadData();
        },
        child: const Icon(Icons.add),
      ),
    );
  }

  Widget _buildExpensesTab() {
    if (_expenses.isEmpty) {
      return const Center(
        child: Text(
          'Nenhuma despesa ainda.',
          style: TextStyle(fontSize: 16, color: Colors.grey),
        ),
      );
    }

    final currencyFormat = NumberFormat.currency(locale: 'pt_BR', symbol: 'R\$');

    return RefreshIndicator(
      onRefresh: _loadData,
      child: ListView.separated(
        padding: const EdgeInsets.all(16),
        itemCount: _expenses.length,
        separatorBuilder: (_, __) => const SizedBox(height: 8),
        itemBuilder: (context, index) {
          final expense = _expenses[index];
          return Card(
            child: ListTile(
              leading: CircleAvatar(
                backgroundColor: Theme.of(context).colorScheme.primary,
                child: const Icon(Icons.receipt, color: Colors.white),
              ),
              title: Text(expense.description),
              subtitle: Text(
                '${expense.expenseDate} • ${expense.paidByName.isNotEmpty ? expense.paidByName : expense.paidBy}',
              ),
              trailing: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Text(
                    currencyFormat.format(expense.amount),
                    style: const TextStyle(
                      fontWeight: FontWeight.bold,
                      fontSize: 14,
                    ),
                  ),
                  IconButton(
                    icon: const Icon(Icons.delete_outline, color: Colors.red),
                    onPressed: () => _deleteExpense(expense.id),
                  ),
                ],
              ),
            ),
          );
        },
      ),
    );
  }

  Widget _buildBalancesTab() {
    if (_balances.isEmpty) {
      return const Center(
        child: Text(
          'Nenhum balanço ainda.',
          style: TextStyle(fontSize: 16, color: Colors.grey),
        ),
      );
    }

    final currencyFormat = NumberFormat.currency(locale: 'pt_BR', symbol: 'R\$');

    return ListView.separated(
      padding: const EdgeInsets.all(16),
      itemCount: _balances.length,
      separatorBuilder: (_, __) => const SizedBox(height: 8),
      itemBuilder: (context, index) {
        final balance = _balances[index];
        final isPositive = balance.balance >= 0;
        return Card(
          child: ListTile(
            leading: CircleAvatar(
              backgroundColor: isPositive ? Colors.green : Colors.red,
              child: Text(
                balance.name.isNotEmpty ? balance.name[0].toUpperCase() : '?',
                style: const TextStyle(color: Colors.white),
              ),
            ),
            title: Text(balance.name),
            trailing: Text(
              currencyFormat.format(balance.balance),
              style: TextStyle(
                fontWeight: FontWeight.bold,
                fontSize: 16,
                color: isPositive ? Colors.green : Colors.red,
              ),
            ),
          ),
        );
      },
    );
  }

  Widget _buildMembersTab() {
    if (_members.isEmpty) {
      return const Center(
        child: Text(
          'Nenhum membro.',
          style: TextStyle(fontSize: 16, color: Colors.grey),
        ),
      );
    }

    return ListView.separated(
      padding: const EdgeInsets.all(16),
      itemCount: _members.length,
      separatorBuilder: (_, __) => const SizedBox(height: 8),
      itemBuilder: (context, index) {
        final member = _members[index];
        return Card(
          child: ListTile(
            leading: CircleAvatar(
              backgroundColor: Theme.of(context).colorScheme.primary,
              child: Text(
                member.name[0].toUpperCase(),
                style: const TextStyle(color: Colors.white),
              ),
            ),
            title: Text(member.name),
            subtitle: Text(member.email),
          ),
        );
      },
    );
  }
}
