import 'dart:async';
import 'package:flutter/material.dart';
import '../service/auth_service.dart';
import '../service/config_service.dart';
import 'groups_screen.dart';
import 'join_group_screen.dart';
import 'register_screen.dart';

class LoginScreen extends StatefulWidget {
  final String? joinGroupId;

  const LoginScreen({super.key, this.joinGroupId});

  @override
  State<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends State<LoginScreen> {
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();
  final _authService = AuthService();
  bool _loading = false;

  int _tapCount = 0;
  Timer? _tapTimer;

  @override
  void dispose() {
    _emailController.dispose();
    _passwordController.dispose();
    _tapTimer?.cancel();
    super.dispose();
  }

  void _onLoginTap() {
    _tapCount++;
    _tapTimer?.cancel();

    if (_tapCount >= 2) {
      _tapCount = 0;
      _showUrlDialog();
      return;
    }

    _tapTimer = Timer(const Duration(milliseconds: 350), () {
      _tapCount = 0;
      if (!_loading) _login();
    });
  }

  Future<void> _login() async {
    if (_emailController.text.isEmpty || _passwordController.text.isEmpty) {
      _showError('Preencha todos os campos');
      return;
    }
    setState(() => _loading = true);
    try {
      await _authService.login(
        _emailController.text.trim(),
        _passwordController.text,
      );
      if (!mounted) return;
      final destination = widget.joinGroupId != null
          ? JoinGroupScreen(groupId: widget.joinGroupId!)
          : const GroupsScreen();
      Navigator.pushReplacement(
        context,
        MaterialPageRoute(builder: (_) => destination),
      );
    } catch (e) {
      _showError(e.toString().replaceFirst('Exception: ', ''));
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  void _showUrlDialog() {
    final urlController = TextEditingController(text: ConfigService.baseUrl);

    showDialog(
      context: context,
      builder: (context) {
        return AlertDialog(
          title: const Text('URL do servidor'),
          content: TextField(
            controller: urlController,
            keyboardType: TextInputType.url,
            autocorrect: false,
            decoration: const InputDecoration(
              labelText: 'Ex: https://xxxx.ngrok.io',
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
                final url = urlController.text.trim();
                if (url.isEmpty) return;
                await ConfigService.save(url);
                if (!context.mounted) return;
                Navigator.pop(context);
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(content: Text('Servidor: ${ConfigService.baseUrl}')),
                );
              },
              child: const Text('Salvar'),
            ),
          ],
        );
      },
    );
  }

  void _showError(String message) {
    if (!mounted) return;
    ScaffoldMessenger.of(context)
        .showSnackBar(SnackBar(content: Text(message)));
  }

  @override
  Widget build(BuildContext context) {
    final screenWidth = MediaQuery.of(context).size.width;
    final formWidth = screenWidth > 600 ? 480.0 : screenWidth;

    return Scaffold(
      body: Center(
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(24),
          child: SizedBox(
            width: formWidth,
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                const Icon(Icons.receipt_long, size: 72, color: Colors.indigo),
                const SizedBox(height: 16),
                Text(
                  'Racha Histórico',
                  textAlign: TextAlign.center,
                  style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                ),
                const SizedBox(height: 40),
                TextField(
                  controller: _emailController,
                  keyboardType: TextInputType.emailAddress,
                  decoration: const InputDecoration(
                    labelText: 'Email',
                    border: OutlineInputBorder(),
                    prefixIcon: Icon(Icons.email_outlined),
                  ),
                ),
                const SizedBox(height: 16),
                TextField(
                  controller: _passwordController,
                  obscureText: true,
                  decoration: const InputDecoration(
                    labelText: 'Senha',
                    border: OutlineInputBorder(),
                    prefixIcon: Icon(Icons.lock_outline),
                  ),
                  onSubmitted: (_) => _onLoginTap(),
                ),
                const SizedBox(height: 24),
                FilledButton(
                  onPressed: _loading ? null : _onLoginTap,
                  child: _loading
                      ? const SizedBox(
                          height: 20,
                          width: 20,
                          child: CircularProgressIndicator(
                              strokeWidth: 2, color: Colors.white),
                        )
                      : const Text('Entrar'),
                ),
                const SizedBox(height: 12),
                TextButton(
                  onPressed: () => Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (_) => RegisterScreen(joinGroupId: widget.joinGroupId),
                    ),
                  ),
                  child: const Text('Não tem conta? Cadastre-se'),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
