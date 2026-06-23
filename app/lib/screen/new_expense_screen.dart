import 'package:flutter/material.dart';
import 'package:image_picker/image_picker.dart';
import '../domain/group.dart';
import '../service/expense_service.dart';
import '../service/notification_service.dart';

class NewExpenseScreen extends StatefulWidget {
  final String groupId;
  final List<GroupMember> members;

  const NewExpenseScreen({
    super.key,
    required this.groupId,
    required this.members,
  });

  @override
  State<NewExpenseScreen> createState() => _NewExpenseScreenState();
}

class _NewExpenseScreenState extends State<NewExpenseScreen> {
  final _formKey = GlobalKey<FormState>();
  final _descriptionController = TextEditingController();
  final _amountController = TextEditingController();
  final ExpenseService _expenseService = ExpenseService();

  DateTime _selectedDate = DateTime.now();
  List<String> _selectedMemberIds = [];
  bool _loading = false;
  bool _processingOCR = false;

  @override
  void initState() {
    super.initState();
    _selectedMemberIds = widget.members.map((m) => m.userId).toList();
  }

  @override
  void dispose() {
    _descriptionController.dispose();
    _amountController.dispose();
    super.dispose();
  }

  Future<void> _pickDate() async {
    final picked = await showDatePicker(
      context: context,
      initialDate: _selectedDate,
      firstDate: DateTime(2020),
      lastDate: DateTime.now(),
    );
    if (picked != null) setState(() => _selectedDate = picked);
  }

  Future<void> _scanReceipt() async {
    final picker = ImagePicker();
    final image = await picker.pickImage(source: ImageSource.camera);
    if (image == null) return;

    setState(() => _processingOCR = true);
    try {
      final bytes = await image.readAsBytes();
      final result = await _expenseService.processOCR(bytes, image.name);
      setState(() {
        if (result.description.isNotEmpty) {
          _descriptionController.text = result.description;
        }
        if (result.amount > 0) {
          _amountController.text = result.amount.toStringAsFixed(2);
        }
      });
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(e.toString().replaceFirst('Exception: ', ''))),
      );
    } finally {
      if (mounted) setState(() => _processingOCR = false);
    }
  }

  Future<void> _save() async {
    if (!_formKey.currentState!.validate()) return;
    if (_selectedMemberIds.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Selecione ao menos um participante')),
      );
      return;
    }

    setState(() => _loading = true);
    try {
      final dateStr =
          '${_selectedDate.year}-${_selectedDate.month.toString().padLeft(2, '0')}-${_selectedDate.day.toString().padLeft(2, '0')}';
      final description = _descriptionController.text.trim();
      final amount = double.parse(_amountController.text.replaceAll(',', '.'));
      await _expenseService.createExpense(
        groupId: widget.groupId,
        description: description,
        amount: amount,
        expenseDate: dateStr,
        splitUserIds: _selectedMemberIds,
      );
      await NotificationService.showExpenseAdded(description, amount);
      if (!mounted) return;
      Navigator.pop(context);
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(e.toString().replaceFirst('Exception: ', ''))),
      );
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final screenWidth = MediaQuery.of(context).size.width;
    final formWidth = screenWidth > 600 ? 600.0 : screenWidth;
    final dateLabel =
        '${_selectedDate.day.toString().padLeft(2, '0')}/${_selectedDate.month.toString().padLeft(2, '0')}/${_selectedDate.year}';

    return Scaffold(
      appBar: AppBar(
        title: const Text('Nova Despesa'),
        actions: [
          _processingOCR
              ? const Padding(
                  padding: EdgeInsets.all(12),
                  child: SizedBox(
                    width: 24,
                    height: 24,
                    child: CircularProgressIndicator(strokeWidth: 2),
                  ),
                )
              : IconButton(
                  icon: const Icon(Icons.camera_alt),
                  onPressed: _scanReceipt,
                  tooltip: 'Escanear comprovante',
                ),
        ],
      ),
      body: Center(
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(16),
          child: SizedBox(
            width: formWidth,
            child: Form(
              key: _formKey,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  TextFormField(
                    controller: _descriptionController,
                    decoration: const InputDecoration(
                      labelText: 'Descrição',
                      border: OutlineInputBorder(),
                    ),
                    validator: (value) {
                      if (value == null || value.isEmpty) return 'Informe a descrição';
                      return null;
                    },
                  ),
                  const SizedBox(height: 16),
                  TextFormField(
                    controller: _amountController,
                    keyboardType:
                        const TextInputType.numberWithOptions(decimal: true),
                    decoration: const InputDecoration(
                      labelText: 'Valor (R\$)',
                      border: OutlineInputBorder(),
                    ),
                    validator: (value) {
                      if (value == null || value.isEmpty) return 'Informe o valor';
                      final parsed =
                          double.tryParse(value.replaceAll(',', '.'));
                      if (parsed == null || parsed <= 0) return 'Valor inválido';
                      return null;
                    },
                  ),
                  const SizedBox(height: 16),
                  OutlinedButton.icon(
                    onPressed: _pickDate,
                    icon: const Icon(Icons.calendar_today),
                    label: Text('Data: $dateLabel'),
                  ),
                  const SizedBox(height: 16),
                  const Text(
                    'Participantes',
                    style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
                  ),
                  const SizedBox(height: 8),
                  ...widget.members.map((member) {
                    final isSelected = _selectedMemberIds.contains(member.userId);
                    return CheckboxListTile(
                      title: Text(member.name),
                      subtitle: Text(member.email),
                      value: isSelected,
                      onChanged: (checked) {
                        setState(() {
                          if (checked == true) {
                            _selectedMemberIds.add(member.userId);
                          } else {
                            _selectedMemberIds.remove(member.userId);
                          }
                        });
                      },
                    );
                  }),
                  const SizedBox(height: 24),
                  FilledButton(
                    onPressed: _loading ? null : _save,
                    style: FilledButton.styleFrom(
                      padding: const EdgeInsets.symmetric(vertical: 16),
                    ),
                    child: _loading
                        ? const SizedBox(
                            height: 20,
                            width: 20,
                            child: CircularProgressIndicator(
                              strokeWidth: 2,
                              color: Colors.white,
                            ),
                          )
                        : const Text('Salvar', style: TextStyle(fontSize: 16)),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
