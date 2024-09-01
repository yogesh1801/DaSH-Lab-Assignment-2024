#ifndef TABLE_CPP
#define TABLE_CPP


#include <vector>
#include <string>
#include <memory>

class Table {
private:
    std::vector<std::vector<std::string>> data;
    std::string name;
    int refCount;
    bool isLockable;

public:
    Table(int rows, int columns, std::string tablename) : refCount(0), isLockable(false), name(tablename) {
        data.resize(rows, std::vector<std::string>(columns));
    }

    void incrementRefCount() {
        ++refCount;
    }

    std::string getName() {
        return name;
    }

    void decrementRefCount() {
        if (refCount > 0) --refCount;
    }

    int getRefCount() const {
        return refCount;
    }

    void setLockable(bool lockable) {
        isLockable = lockable;
    }

    bool getLockable() const {
        return isLockable;
    }

    void setData(int row, int col, const std::string& value) {
        if (row < data.size() && col < data[0].size()) {
            data[row][col] = value;
        }
    }

    void setRow(int row, const std::vector<std::string>& rowData) {
        if (row < data.size()) {
            if (rowData.size() == data[row].size()) {
                data[row] = rowData;
            } else {
                throw std::invalid_argument("Row data size does not match table column count");
            }
        } else {
            throw std::out_of_range("Row index out of range");
        }
    }

    std::string getData(int row, int col) const {
        if (row < data.size() && col < data[0].size()) {
            return data[row][col];
        }
        return "";
    }

    int getRowCount() const {
        return data.size();
    }

    int getColumnCount() const {
        return data.empty() ? 0 : data[0].size();
    }
};

#endif