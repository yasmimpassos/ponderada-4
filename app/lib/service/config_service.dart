import 'package:shared_preferences/shared_preferences.dart';

class ConfigService {
  static const String _urlKey = 'backend_url';
  static const String _defaultUrl = 'http://10.0.2.2:8080';

  static String _baseUrl = _defaultUrl;

  static String get baseUrl => _baseUrl;

  static Future<void> load() async {
    final prefs = await SharedPreferences.getInstance();
    _baseUrl = prefs.getString(_urlKey) ?? _defaultUrl;
  }

  static Future<void> save(String url) async {
    var cleaned = url.trim();
    if (cleaned.endsWith('/')) cleaned = cleaned.substring(0, cleaned.length - 1);
    if (!cleaned.startsWith('http://') && !cleaned.startsWith('https://')) {
      cleaned = 'http://$cleaned';
    }
    _baseUrl = cleaned;
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_urlKey, _baseUrl);
  }
}
