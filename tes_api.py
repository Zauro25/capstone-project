import pytest
import requests
import json
from datetime import datetime

BASE_URL = "http://localhost:8080/api/v1"
TEST_OUTPUT = []

def log_test_result(test_name, status, response=None):
    """Helper function to log test results for PDF output"""
    result = {
        "test_name": test_name,
        "status": status,
        "response": response.get("message") if response and isinstance(response, dict) else str(response),
        "timestamp": datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    }
    TEST_OUTPUT.append(result)

@pytest.fixture
def auth_token():
    """Register and login a user to get JWT token"""
    # Register a test user
    register_payload = {
        "username": "testuser",
        "password": "testpassword123",
        "role": "AdminPerpustakaan"
    }
    register_response = requests.post(f"{BASE_URL}/auth/register", json=register_payload)
    log_test_result("Register User", "PASSED" if register_response.status_code == 201 else "FAILED", register_response.json())
    assert register_response.status_code == 201, f"Register failed: {register_response.text}"

    # Login to get token
    login_payload = {
        "username": "testuser",
        "password": "testpassword123"
    }
    login_response = requests.post(f"{BASE_URL}/auth/login", json=login_payload)
    log_test_result("Login User", "PASSED" if login_response.status_code == 200 else "FAILED", login_response.json())
    assert login_response.status_code == 200, f"Login failed: {login_response.text}"

    token = login_response.json().get("data", {}).get("token")
    assert token, "Token not found in login response"
    return token

def test_register_user_invalid_role():
    """Test register with invalid role"""
    payload = {
        "username": "invaliduser",
        "password": "testpassword123",
        "role": "InvalidRole"
    }
    response = requests.post(f"{BASE_URL}/auth/register", json=payload)
    log_test_result("Register User (Invalid Role)", "PASSED" if response.status_code == 400 else "FAILED", response.json())
    assert response.status_code == 400
    assert "error" in response.json()

def test_login_invalid_credentials():
    """Test login with wrong password"""
    payload = {
        "username": "testuser",
        "password": "wrongpassword"
    }
    response = requests.post(f"{BASE_URL}/auth/login", json=payload)
    log_test_result("Login User (Invalid Credentials)", "PASSED" if response.status_code == 401 else "FAILED", response.json())
    assert response.status_code == 401
    assert "error" in response.json()

def test_create_perpustakaan(auth_token):
    """Test creating a new perpustakaan"""
    headers = {"Authorization": f"Bearer {auth_token}"}
    payload = {
        "nama": "Perpustakaan Test",
        "alamat": "Jl. Test No. 123",
        "jenis": "Umum",
        "nomor_induk": "TEST123",
        "jumlah_sdm": 5,
        "jumlah_pengunjung": 100,
        "jumlah_anggota": 50,
        "status_verifikasi": "Belum Diverifikasi"
    }
    response = requests.post(f"{BASE_URL}/admin-perpustakaan/perpustakaan", json=payload, headers=headers)
    log_test_result("Create Perpustakaan", "PASSED" if response.status_code == 201 else "FAILED", response.json())
    assert response.status_code == 201
    assert "data" in response.json()
    return response.json().get("data", {}).get("id")

def test_create_perpustakaan_invalid_input(auth_token):
    """Test creating perpustakaan with missing required fields"""
    headers = {"Authorization": f"Bearer {auth_token}"}
    payload = {
        "alamat": "Jl. Test No. 123",
        "jumlah_sdm": 5
    }
    response = requests.post(f"{BASE_URL}/admin-perpustakaan/perpustakaan", json=payload, headers=headers)
    log_test_result("Create Perpustakaan (Invalid Input)", "PASSED" if response.status_code == 400 else "FAILED", response.json())
    assert response.status_code == 400
    assert "error" in response.json()

def test_get_perpustakaan_by_id(auth_token, perpustakaan_id):
    """Test getting perpustakaan by ID"""
    headers = {"Authorization": f"Bearer {auth_token}"}
    response = requests.get(f"{BASE_URL}/admin-perpustakaan/perpustakaan/{perpustakaan_id}", headers=headers)
    log_test_result("Get Perpustakaan by ID", "PASSED" if response.status_code == 200 else "FAILED", response.json())
    assert response.status_code == 200
    assert "data" in response.json()

def test_submit_verifikasi(auth_token, perpustakaan_id):
    """Test submitting perpustakaan for verification"""
    headers = {"Authorization": f"Bearer {auth_token}"}
    response = requests.post(f"{BASE_URL}/admin-perpustakaan/perpustakaan/{perpustakaan_id}/submit-verifikasi", headers=headers)
    log_test_result("Submit Perpustakaan for Verification", "PASSED" if response.status_code == 200 else "FAILED", response.json())
    assert response.status_code == 200
    assert "message" in response.json()

def test_add_sdm(auth_token, perpustakaan_id):
    """Test adding SDM to perpustakaan"""
    headers = {"Authorization": f"Bearer {auth_token}"}
    payload = {
        "nama": "John Doe",
        "jabatan": "Pustakawan",
        "pendidikan_terakhir": "S1",
        "status_kepegawaian": "Tetap"
    }
    response = requests.post(f"{BASE_URL}/admin-perpustakaan/perpustakaan/{perpustakaan_id}/sdm", json=payload, headers=headers)
    log_test_result("Add SDM", "PASSED" if response.status_code == 201 else "FAILED", response.json())
    assert response.status_code == 201
    assert "data" in response.json()

def test_add_anggota(auth_token, perpustakaan_id):
    """Test adding anggota to perpustakaan"""
    headers = {"Authorization": f"Bearer {auth_token}"}
    payload = {
        "nama": "Jane Doe",
        "tanggal_daftar": "2025-06-25",
        "status_aktif": True,
        "jenis_kelamin": "Perempuan",
        "pekerjaan": "Mahasiswa"
    }
    response = requests.post(f"{BASE_URL}/admin-perpustakaan/perpustakaan/{perpustakaan_id}/anggota", json=payload, headers=headers)
    log_test_result("Add Anggota", "PASSED" if response.status_code == 201 else "FAILED", response.json())
    assert response.status_code == 201
    assert "data" in response.json()

def test_add_pengunjung(auth_token, perpustakaan_id):
    """Test adding pengunjung to perpustakaan"""
    headers = {"Authorization": f"Bearer {auth_token}"}
    payload = {
        "tanggal_kunjungan": "2025-06-25",
        "jumlah_pengunjung": 10,
        "keterangan": "Kunjungan kelompok"
    }
    response = requests.post(f"{BASE_URL}/admin-perpustakaan/perpustakaan/{perpustakaan_id}/pengunjung", json=payload, headers=headers)
    log_test_result("Add Pengunjung", "PASSED" if response.status_code == 201 else "FAILED", response.json())
    assert response.status_code == 201
    assert "data" in response.json()

def test_unauthorized_access():
    """Test accessing protected endpoint without token"""
    response = requests.get(f"{BASE_URL}/admin-perpustakaan/perpustakaan/1")
    log_test_result("Unauthorized Access", "PASSED" if response.status_code == 401 else "FAILED", response.json())
    assert response.status_code == 401
    assert "error" in response.json()

if __name__ == "__main__":
    pytest.main(["-v", __file__])
    # Save test results to JSON for PDF generation
    with open("test_results.json", "w") as f:
        json.dump(TEST_OUTPUT, f, indent=2)