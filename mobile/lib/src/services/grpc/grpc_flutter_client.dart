import 'dart:typed_data';

import 'package:flutter/foundation.dart';
import 'package:grpc/grpc.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:flutter/services.dart' show ByteData, rootBundle;

class GrpcFlutterClient {
  static ClientChannel? _clientChannel;

  static Future<ClientChannel> getClient() async {
    ByteData byteData = await rootBundle.load(Constants.caCertPath);
    Uint8List bytes = byteData.buffer.asUint8List();

    String myHost;
    if (defaultTargetPlatform == TargetPlatform.iOS) {
      myHost = 'localhost';
    } else if (defaultTargetPlatform == TargetPlatform.android) {
      myHost = '10.0.2.2';
    } else {
      myHost = '';
    }

    if (_clientChannel == null) {
      _clientChannel = ClientChannel(
        myHost,
        port: 80,
        options: ChannelOptions(
          credentials: ChannelCredentials.secure(
            certificates: bytes,
            onBadCertificate: (certificate, host) => host == '${myHost}:80',
          ),
          idleTimeout: Duration(minutes: 1),
        ),
      );
    }
    return _clientChannel!;
  }
}
