import 'package:flutter_test/flutter_test.dart';

import 'package:bhaavyaad_app/main.dart';

void main() {
  test('lastPricePerKey keeps the newest price per supplier+item', () {
    final ps = [
      Purchase('ACME', 'soap', 22), // newest
      Purchase('ACME', 'soap', 20), // older, ignored
      Purchase('BEST', 'soap', 25),
    ];
    final m = lastPricePerKey(ps);
    expect(m['ACME · soap'], 22);
    expect(m['BEST · soap'], 25);
  });

  testWidgets('shows the price-memory header', (tester) async {
    await tester.pumpWidget(const BhaavyaadApp());
    expect(find.text('Last price you paid'), findsOneWidget);
  });
}
