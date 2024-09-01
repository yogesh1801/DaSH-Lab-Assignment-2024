#include <iostream>
#include <unordered_map>
#include <list>
#include <algorithm>
#include "database.cpp"
#include "transaction.cpp"

Table* createDummyTable(const std::string& name, int rows, int cols) {
    Table* table = new Table(rows, cols, name);
    for (int i = 0; i < rows; ++i) {
        for (int j = 0; j < cols; ++j) {
            table->setData(i, j, "Data_" + std::to_string(i) + "_" + std::to_string(j));
        }
    }
    return table;
}

// Function to check if a transaction can be scheduled
bool canScheduleTransaction(const Transaction& transaction, 
                            const std::unordered_map<std::string, std::list<Transaction*>>& lockMap) {
    if (transaction.tablelock) {
        // For table lock, check if any locks exist for this table
        return lockMap.find(transaction.table->getName()) == lockMap.end() || 
               lockMap.at(transaction.table->getName()).empty();
    }

    // For other locks, check each lock hash
    for (const auto& hash : transaction.lockhashes) {
        if (lockMap.find(hash) != lockMap.end() && !lockMap.at(hash).empty()) {
            return false;
        }
    }
    return true;
}

void scheduleTransaction(Transaction* transaction, 
                         std::unordered_map<std::string, std::list<Transaction*>>& lockMap) {
    if (transaction->tablelock) {
        lockMap[transaction->table->getName()].push_back(transaction);
    } else {
        for (const auto& hash : transaction->lockhashes) {
            lockMap[hash].push_back(transaction);
        }
    }
    std::cout << "Scheduled transaction " << transaction->t_id << std::endl;
}

int main() {
    Database db;
    db.addTable(createDummyTable("Table1", 5, 5));
    db.addTable(createDummyTable("Table2", 4, 4));
    db.addTable(createDummyTable("Table3", 3, 3));

    // Create dummy transactions
    std::vector<Transaction> transactions;
    transactions.push_back(Transaction("T1", db.getTable("Table1"), true, "ACTIVE")); // Table lock
    transactions.push_back(Transaction("T2", db.getTable("Table1"), {Lock("ROW", "Table1", 0)}, "ACTIVE")); // Row lock
    transactions.push_back(Transaction("T3", db.getTable("Table1"), {Lock("COLUMN", "Table1", -1, 2)}, "ACTIVE")); // Column lock
    transactions.push_back(Transaction("T4", db.getTable("Table2"), {Lock("CELL", "Table2", 1, 1)}, "ACTIVE")); // Cell lock
    transactions.push_back(Transaction("T1", db.getTable("Table3"), true, "ACTIVE")); // Another table lock

    // Create lock map
    std::unordered_map<std::string, std::list<Transaction*>> lockMap;

    // Try to schedule transactions
    for (auto& transaction : transactions) {
        if (canScheduleTransaction(transaction, lockMap)) {
            scheduleTransaction(&transaction, lockMap);
        } else {
            std::cout << "Cannot schedule transaction " << transaction.t_id << " due to conflicts" << std::endl;
        }
    }

    std::cout << "\nFinal Lock Status:" << std::endl;
    for (const auto& entry : lockMap) {
        std::cout << "Lock hash: " << entry.first << " - Transactions: ";
        for (const auto& trans : entry.second) {
            std::cout << trans->t_id << " ";
        }
        std::cout << std::endl;
    }

    return 0;
}