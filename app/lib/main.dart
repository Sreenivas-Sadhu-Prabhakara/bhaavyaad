import 'package:flutter/material.dart';

void main() => runApp(const BhaavyaadApp());

/// Bhaavyaad — your purchase-price memory (not live market rates). Records what
/// you paid each supplier per item and recalls the last price before you agree
/// to a new quote. Mirrors the Go journal service.
class BhaavyaadApp extends StatelessWidget {
  const BhaavyaadApp({super.key});
  @override
  Widget build(BuildContext context) => MaterialApp(
        title: 'Bhaavyaad',
        debugShowCheckedModeBanner: false,
        theme: ThemeData(colorSchemeSeed: const Color(0xFF3A6E4F), useMaterial3: true),
        home: const HomePage(),
      );
}

class Purchase {
  final String supplier, item;
  final double price;
  Purchase(this.supplier, this.item, this.price);
}

/// lastPricePerKey returns the latest price per supplier+item from newest-first data.
Map<String, double> lastPricePerKey(List<Purchase> ps) {
  final out = <String, double>{};
  for (final p in ps) {
    final k = '${p.supplier} · ${p.item}';
    out.putIfAbsent(k, () => p.price);
  }
  return out;
}

class HomePage extends StatefulWidget {
  const HomePage({super.key});
  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  final _purchases = <Purchase>[];
  final _sup = TextEditingController();
  final _item = TextEditingController();
  final _price = TextEditingController();

  void _add() {
    final price = double.tryParse(_price.text.trim()) ?? 0;
    if (_sup.text.trim().isEmpty || _item.text.trim().isEmpty || price <= 0) return;
    setState(() {
      _purchases.insert(0, Purchase(_sup.text.trim(), _item.text.trim(), price));
      _price.clear();
    });
  }

  @override
  Widget build(BuildContext context) {
    final last = lastPricePerKey(_purchases);
    return Scaffold(
      appBar: AppBar(
        title: const Text('Bhaavyaad · your price memory'),
        backgroundColor: Theme.of(context).colorScheme.primaryContainer,
      ),
      body: Column(children: [
        Padding(
          padding: const EdgeInsets.all(12),
          child: Row(children: [
            Expanded(child: TextField(controller: _sup, decoration: const InputDecoration(labelText: 'Supplier', border: OutlineInputBorder()))),
            const SizedBox(width: 8),
            Expanded(child: TextField(controller: _item, decoration: const InputDecoration(labelText: 'Item', border: OutlineInputBorder()))),
            const SizedBox(width: 8),
            SizedBox(width: 90, child: TextField(controller: _price, keyboardType: TextInputType.number, decoration: const InputDecoration(labelText: '₹', border: OutlineInputBorder()))),
            const SizedBox(width: 8),
            FilledButton(onPressed: _add, child: const Text('Save')),
          ]),
        ),
        const Padding(
          padding: EdgeInsets.symmetric(horizontal: 12),
          child: Align(alignment: Alignment.centerLeft,
            child: Text('Last price you paid', style: TextStyle(fontWeight: FontWeight.w600))),
        ),
        Expanded(
          child: ListView(children: [
            for (final e in last.entries)
              ListTile(
                title: Text(e.key),
                trailing: Text('₹${e.value.toStringAsFixed(2)}',
                    style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
              ),
          ]),
        ),
        Container(
          width: double.infinity,
          color: Theme.of(context).colorScheme.surfaceContainerHighest,
          padding: const EdgeInsets.all(10),
          child: const Text('This is YOUR purchase memory — not live market rates.',
              textAlign: TextAlign.center, style: TextStyle(fontSize: 12)),
        ),
      ]),
    );
  }
}
