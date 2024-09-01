#ifndef TRANSACTION_CPP
#define TRANSACTION_CPP


#include <vector>
#include <string>
#include <memory>
#include "locks.cpp"
#include "table.cpp"

class Transaction {
public:
    std::string t_id;
    bool tablelock;
    Table* table;
    std::vector<Lock> locks;
    std::string status;
    std::vector<std::string> lockhashes;

    Transaction(const std::string& trans_id, Table* table_ptr, bool table_lock, const std::string& trans_status)
        : t_id(trans_id), tablelock(table_lock), table(table_ptr), status(trans_status)
    {
        if (!tablelock) {
            locks = std::vector<Lock>();
            lockhashes = std::vector<std::string>();
        }
    }

    Transaction(const std::string& trans_id, Table* table_ptr, const std::vector<Lock>& locks_arr, const std::string& trans_status)
        : t_id(trans_id), tablelock(false), table(table_ptr), locks(locks_arr), status(trans_status)
    {
        lockhashes = calculateLockHashes();
    }

private:
    std::string hashFunction(const std::string& base_string) {
        return std::to_string(std::hash<std::string>{}(base_string));
    }

    std::vector<std::string> calculateLockHashes() {
        std::vector<std::string> hashes;

        for (const Lock& lock : locks) {
            if (lock.type == "ROW") {
                std::string rowHash = hashFunction(lock.table_name + "_ROW_" + std::to_string(lock.row));
                hashes.push_back(rowHash);

                int columnCount = table->getColumnCount();
                for (int col = 0; col < columnCount; ++col) {
                    std::string columnHash = hashFunction(lock.table_name + "_COL_" + std::to_string(col));
                    hashes.push_back(columnHash);
                }
            } 
            else if (lock.type == "COLUMN") {
    
                std::string columnHash = hashFunction(lock.table_name + "_COL_" + std::to_string(lock.column));
                hashes.push_back(columnHash);

                int rowCount = table->getRowCount();
                for (int row = 0; row < rowCount; ++row) {
                    std::string columnHash = hashFunction(lock.table_name + "_ROW_" + std::to_string(row));
                    hashes.push_back(columnHash);
                }
            } 
            else if (lock.type == "CELL") {
                std::string rowHash = hashFunction(lock.table_name + "_ROW_" + std::to_string(lock.row));
                hashes.push_back(rowHash);
                std::string columnHash = hashFunction(lock.table_name + "_COL_" + std::to_string(lock.column));
                hashes.push_back(columnHash);
            }
        }

        return hashes;
    }
};

#endif