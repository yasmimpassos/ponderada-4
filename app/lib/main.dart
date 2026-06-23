import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'screen/login_screen.dart';
import 'screen/groups_screen.dart';
import 'screen/join_group_screen.dart';
import 'service/auth_service.dart';
import 'service/config_service.dart';
import 'service/notification_service.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await ConfigService.load();
  await NotificationService.init();
  runApp(const RachaHistoricoApp());
}

class RachaHistoricoApp extends StatelessWidget {
  const RachaHistoricoApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Racha Histórico',
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.indigo),
        useMaterial3: true,
      ),
      home: const SplashScreen(),
    );
  }
}

class SplashScreen extends StatefulWidget {
  const SplashScreen({super.key});

  @override
  State<SplashScreen> createState() => _SplashScreenState();
}

class _SplashScreenState extends State<SplashScreen> {
  final AuthService _authService = AuthService();

  @override
  void initState() {
    super.initState();
    _checkLogin();
  }

  Future<void> _checkLogin() async {
    final joinGroupId = _getJoinGroupFromUrl();
    final token = await _authService.getToken();
    if (!mounted) return;

    if (token != null) {
      if (joinGroupId != null) {
        Navigator.pushReplacement(
          context,
          MaterialPageRoute(
            builder: (_) => JoinGroupScreen(groupId: joinGroupId),
          ),
        );
      } else {
        Navigator.pushReplacement(
          context,
          MaterialPageRoute(builder: (_) => const GroupsScreen()),
        );
      }
    } else {
      Navigator.pushReplacement(
        context,
        MaterialPageRoute(
          builder: (_) => LoginScreen(joinGroupId: joinGroupId),
        ),
      );
    }
  }

  String? _getJoinGroupFromUrl() {
    if (!kIsWeb) return null;
    final uri = Uri.base;
    return uri.queryParameters['join'];
  }

  @override
  Widget build(BuildContext context) {
    return const Scaffold(
      body: Center(child: CircularProgressIndicator()),
    );
  }
}
